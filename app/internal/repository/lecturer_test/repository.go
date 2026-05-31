package lecturertest

import (
	"context"
	"database/sql"
	"strings"
	"testum-engine/app/internal/adapter/db"
	"time"

	"go.uber.org/zap"
)

type Repository interface {
	Create(ctx context.Context, lecturerID int, test Test) (int, error)
	Delete(ctx context.Context, testID int) (bool, error)
	GetByID(ctx context.Context, testID int) (TestInfo, error)
	GetByLecturer(ctx context.Context, lecturerID int) ([]TestInfo, error)
	GetGroups(ctx context.Context, testID int, year int) ([]GroupInfo, error)
}

type repository struct {
	db  db.Executor
	log *zap.Logger
	now func() time.Time
}

func NewRepository(db db.Executor, log *zap.Logger) Repository {
	return &repository{
		db:  db,
		log: log,
		now: time.Now,
	}
}

// ==================== CREATE ====================
func (r *repository) Create(ctx context.Context, lecturerID int, test Test) (int, error) {
	// 1. test
	res, err := r.db.ExecContext(ctx, `
		insert into tests (owner_id, name, file_name, date_created)
		values (?, ?, ?, CURRENT_TIMESTAMP)
	`, lecturerID, test.Name, test.FileName)

	if err != nil {
		r.log.Error("failed to insert test", zap.Error(err))
		return 0, ErrCreateFailed
	}

	testID, err := res.LastInsertId()
	if err != nil {
		r.log.Error("failed to get test id", zap.Error(err))
		return 0, err
	}

	// 2. tasks + answers
	for _, task := range test.Tasks {
		taskRes, err := r.db.ExecContext(ctx, `
			insert into tasks (test_id, text, image_url, is_hard)
			values (?, ?, ?, ?)
		`, testID, task.Text, task.ImageURL, task.IsHard)

		if err != nil {
			r.log.Error("failed to insert task", zap.Error(err))
			return 0, ErrCreateFailed
		}

		taskID, err := taskRes.LastInsertId()
		if err != nil {
			r.log.Error("failed to get task id", zap.Error(err))
			return 0, err
		}

		for _, ans := range task.Answers {
			_, err := r.db.ExecContext(ctx, `
				insert into answers (task_id, text, image_url, is_correct)
				values (?, ?, ?, ?)
			`, taskID, ans.Text, ans.ImageURL, ans.IsCorrect)

			if err != nil {
				r.log.Error("failed to insert answer", zap.Error(err))
				return 0, ErrCreateFailed
			}
		}
	}

	// 3. permissions
	for _, group := range test.Groups {
		_, err := r.db.ExecContext(ctx, `
			insert into test_permissions (test_id, `+"`group`"+`)
			values (?, ?)
		`, testID, group)

		if err != nil {
			r.log.Error("failed to insert permission", zap.Error(err))
			return 0, ErrCreateFailed
		}
	}

	return int(testID), nil
}

// ==================== DELETE ====================
func (r *repository) Delete(ctx context.Context, testID int) (bool, error) {
	res, err := r.db.ExecContext(ctx,
		`delete from tests where id = ?`, testID,
	)
	if err != nil {
		r.log.Error("failed to delete test",
			zap.Error(err),
			zap.Int("test_id", testID),
		)
		return false, ErrDeleteFailed
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, ErrTestNotFound
	}

	return true, nil
}

// ==================== GET BY ID ====================
func (r *repository) GetByID(ctx context.Context, testID int) (TestInfo, error) {
	query := `
		select
			t.id,
			t.name,
			(select count(*) from tasks where test_id = t.id),
			(select count(*) from tasks where test_id = t.id and is_hard = 1),
			t.file_name,
			t.date_created
		from tests t
		where t.id = ?
	`

	var result TestInfo

	err := r.db.QueryRowContext(ctx, query, testID).Scan(
		&result.ID,
		&result.Name,
		&result.CntQuestions,
		&result.CntHardQuestions,
		&result.FileName,
		&result.DateCreated,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return TestInfo{}, ErrTestNotFound
		}
		r.log.Error("failed to get test", zap.Error(err))
		return TestInfo{}, err
	}

	// groups
	rows, err := r.db.QueryContext(ctx, `
		select `+"`group`"+` from test_permissions where test_id = ?
	`, testID)
	if err != nil {
		return TestInfo{}, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var g string
		if err := rows.Scan(&g); err != nil {
			return TestInfo{}, err
		}
		result.Groups = append(result.Groups, g)
	}

	if err = rows.Err(); err != nil {
		return TestInfo{}, err
	}

	return result, nil
}

// ==================== GET BY LECTURER ====================
func (r *repository) GetByLecturer(ctx context.Context, lecturerID int) ([]TestInfo, error) {
	query := `
		select
			t.id,
			t.name,
			(select count(*) from tasks where test_id = t.id),
			(select count(*) from tasks where test_id = t.id and is_hard = 1),
			t.file_name,
			t.date_created
		from tests t
		where t.owner_id = ?
		order by t.date_created desc
	`

	rows, err := r.db.QueryContext(ctx, query, lecturerID)
	if err != nil {
		r.log.Error("failed to get tests", zap.Error(err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []TestInfo
	var ids []int

	for rows.Next() {
		var item TestInfo

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.CntQuestions,
			&item.CntHardQuestions,
			&item.FileName,
			&item.DateCreated,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
		ids = append(ids, item.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return result, nil
	}

	in := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		in[i] = "?"
		args[i] = id
	}

	queryGroups := `
		select test_id, ` + "`group`" + `
		from test_permissions
		where test_id in (` + strings.Join(in, ",") + `)
	`

	groupRows, err := r.db.QueryContext(ctx, queryGroups, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = groupRows.Close() }()

	groupMap := make(map[int][]string)

	for groupRows.Next() {
		var (
			testID int
			group  string
		)

		if err := groupRows.Scan(&testID, &group); err != nil {
			return nil, err
		}

		groupMap[testID] = append(groupMap[testID], group)
	}

	for i := range result {
		result[i].Groups = groupMap[result[i].ID]
	}

	return result, nil
}

// ==================== GET GROUPS ====================
func (r *repository) GetGroups(ctx context.Context, testID int, year int) ([]GroupInfo, error) {
	start, end := r.getAcademicRange(year)

	var (
		rows *sql.Rows
		err  error
	)

	if year == 0 {
		query := `
			SELECT 
				res.group_name,
				COALESCE(
					MAX(CASE WHEN res.source = 'st' THEN res.members_count END),
					MAX(CASE WHEN res.source = 'tp' THEN res.members_count END)
				) as members_count
			FROM (
				SELECT 
					tp.` + "`group`" + ` AS group_name,
					COUNT(u.id) as members_count,
					'tp' as source
				FROM test_permissions tp
				LEFT JOIN users u 
					ON u.` + "`group`" + ` = tp.` + "`group`" + `
					AND u.date_modified >= ?
					AND u.date_modified < ?
				WHERE tp.test_id = ?
				GROUP BY tp.` + "`group`" + `

				UNION ALL

				SELECT 
					st.` + "`group`" + ` AS group_name,
					COUNT(DISTINCT st.student_id) as members_count,
					'st' as source
				FROM student_tests st
				WHERE st.test_id = ?
				  AND st.date_start >= ?
				  AND st.date_start < ?
				GROUP BY st.` + "`group`" + `
			) res
			GROUP BY res.group_name
		`

		rows, err = r.db.QueryContext(ctx, query, start, end, testID, testID, start, end)
	} else {
		query := `
			SELECT 
				st.` + "`group`" + `,
				COUNT(DISTINCT st.student_id)
			FROM student_tests st
			WHERE st.test_id = ?
			  AND st.date_start >= ?
			  AND st.date_start < ?
			GROUP BY st.` + "`group`" + `
		`

		rows, err = r.db.QueryContext(ctx, query, testID, start, end)
	}

	if err != nil {
		r.log.Error("GetGroups query failed",
			zap.Error(err),
			zap.Int("test_id", testID),
			zap.Int("year", year),
		)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []GroupInfo

	for rows.Next() {
		var g GroupInfo
		if err := rows.Scan(&g.GroupName, &g.MembersCount); err != nil {
			return nil, err
		}
		result = append(result, g)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, ErrTestNotFound
	}

	return result, nil
}

// ==================== TIME ====================
func (r *repository) getAcademicRange(year int) (time.Time, time.Time) {
	now := r.now()

	currentYear := now.Year()
	august := time.Date(currentYear, time.August, 1, 0, 0, 0, 0, now.Location())

	if now.Before(august) {
		start := august.AddDate(-1-year, 0, 0)
		end := august.AddDate(-year, 0, 0)
		return start, end
	}

	start := august.AddDate(-year, 0, 0)
	end := august.AddDate(1-year, 0, 0)
	return start, end
}

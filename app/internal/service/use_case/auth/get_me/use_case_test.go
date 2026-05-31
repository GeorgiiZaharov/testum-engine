package getme

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	ldapapdap "testum-engine/app/internal/adapter/ldap"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

//
// =========================
// MOCK LDAP
// =========================
//

type mockLDAP struct {
	GetInfoFunc func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error)
}

func (m *mockLDAP) GetInfo(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
	return m.GetInfoFunc(ctx, login)
}

//
// =========================
// TEST ENV
// =========================
//

type testEnv struct {
	uc  *UseCase
	ctx context.Context
}

func setup(t *testing.T) *testEnv {
	t.Helper()

	ctx := context.Background()

	database, cleanup, err := bootstrap.Setup(bootstrap.Config{
		DBOptions: db.DBOptions{
			Path: ":memory:",
		},
		Migrations: "../../../../../migrations/",
	})
	require.NoError(t, err)
	t.Cleanup(cleanup)

	fx := fixtures.New(database)
	require.NoError(t, fx.Reset(ctx))
	require.NoError(t, fx.SeedAll(ctx))

	return &testEnv{
		uc: NewUseCase(
			NewFactory(database, zap.NewNop()),
			&mockLDAP{},
			zap.NewNop(),
		),
		ctx: ctx,
	}
}

//
// =========================
// TESTS
// =========================
//

func TestSuccess_WithFixtures(t *testing.T) {
	env := setup(t)

	group := "A-101"
	ldap := &mockLDAP{
		GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
			return &ldapapdap.LdapUserInfo{
				Login: "student1",
				Name:  "Student One",
				Mail:  "student1@mail.com",
				Group: &group,
			}, nil
		},
	}

	env.uc = NewUseCase(
		NewFactory(env.uc.factory.(*factory).db, zap.NewNop()),
		ldap,
		zap.NewNop(),
	)

	// Выполним запрос
	resp, err := env.uc.Execute(env.ctx, GetMeRequest{
		UserID: 1,
	})

	// Проверки
	require.NoError(t, err)
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, "student1", resp.Login)
	assert.Equal(t, "student1@mail.com", resp.Mail)
	assert.Equal(t, "Student One", resp.Name)
	assert.Equal(t, "A-101", *resp.Group)
}

func TestInvalidInput(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, GetMeRequest{
		UserID: 0,
	})

	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestLDAPGetInfoFail(t *testing.T) {
	env := setup(t)

	env.uc = NewUseCase(
		env.uc.factory,
		&mockLDAP{
			GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
				return nil, errors.New("ldap down")
			},
		},
		zap.NewNop(),
	)

	_, err := env.uc.Execute(env.ctx, GetMeRequest{
		UserID: 4,
	})

	assert.ErrorIs(t, err, ErrLDAPFailed)
}

func TestUserNotFound(t *testing.T) {
	env := setup(t)

	env.uc = NewUseCase(
		env.uc.factory,
		&mockLDAP{
			GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
				return nil, errors.New("ldap down")
			},
		},
		zap.NewNop(),
	)

	_, err := env.uc.Execute(env.ctx, GetMeRequest{
		UserID: 9999,
	})
	fmt.Println(err)

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestRefreshDataFromLDAP(t *testing.T) {
	env := setup(t)

	ldap := &mockLDAP{
		GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
			return &ldapapdap.LdapUserInfo{
				Login: "old_student",
				Name:  "Old Student",
				Mail:  "old_studentNEW@mail.com",
				Group: ptrToString("B-202"),
			}, nil
		},
	}

	env.uc = NewUseCase(
		NewFactory(env.uc.factory.(*factory).db, zap.NewNop()),
		ldap,
		zap.NewNop(),
	)

	// Выполним запрос
	resp, err := env.uc.Execute(env.ctx, GetMeRequest{
		UserID: 4,
	})

	// Проверки
	require.NoError(t, err)
	assert.Equal(t, 4, resp.ID)
	assert.Equal(t, "old_student", resp.Login)
	assert.Equal(t, "old_studentNEW@mail.com", resp.Mail)
	assert.Equal(t, "Old Student", resp.Name)
	assert.Equal(t, "B-202", *resp.Group) // Ожидаем обновленную группу
}

func ptrToString(s string) *string {
	return &s
}

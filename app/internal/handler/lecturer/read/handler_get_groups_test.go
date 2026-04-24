package read

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testum-engine/app/internal/handler/lecturer/read/dto"
	"testum-engine/app/internal/handler/middleware"
	getgroups "testum-engine/app/internal/service/use_case/lecturer/get_groups"
)

type getGroupsUCStub struct {
	executeFn func(ctx context.Context, req getgroups.GetGroupsRequest) (getgroups.GetGroupsResponse, error)
}

func (s *getGroupsUCStub) Execute(
	ctx context.Context,
	req getgroups.GetGroupsRequest,
) (getgroups.GetGroupsResponse, error) {
	return s.executeFn(ctx, req)
}

func TestGetGroups_Success(t *testing.T) {
	uc := &getGroupsUCStub{
		executeFn: func(
			ctx context.Context,
			req getgroups.GetGroupsRequest,
		) (getgroups.GetGroupsResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 5, req.TestID)
			require.Equal(t, 2026, req.Year)

			return getgroups.GetGroupsResponse{
				Groups: []getgroups.GroupInfo{
					{
						GroupName:    "A-01",
						MembersCount: 25,
					},
				},
			}, nil
		},
	}

	h := New(nil, nil, uc, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests/5/groups?year=2026",
		nil,
	)
	req.SetPathValue("test_id", "5")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetGroups(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.GetGroupsResponse

	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Len(t, resp.Groups, 1)
}

func TestGetGroups_MissingYear(t *testing.T) {
	h := New(nil, nil, &getGroupsUCStub{}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests/5/groups",
		nil,
	)
	req.SetPathValue("test_id", "5")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetGroups(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

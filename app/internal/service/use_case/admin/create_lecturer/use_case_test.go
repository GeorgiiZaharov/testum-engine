package createlecturer

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
	userrepo "testum-engine/app/internal/repository/user"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

// =========================
// MOCK LDAP
// =========================

type mockLdap struct {
	GetInfoFunc func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error)
}

func (m *mockLdap) GetInfo(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
	return m.GetInfoFunc(ctx, login)
}

// =========================
// TEST ENV
// =========================

type testEnv struct {
	uc   *UseCase
	ctx  context.Context
	repo userrepo.Repository
}

func setup(t *testing.T) *testEnv {
	t.Helper()

	ctx := context.Background()

	database, cleanup, err := bootstrap.Setup(bootstrap.Config{
		DBOptions: db.DBOptions{
			Host: "localhost",
			Port: "3306",
			User: "testum_user",
			Pass: "testum_pass",
			Name: "testum",
		},
		Migrations: "../../../../../migrations/",
	})
	require.NoError(t, err)
	t.Cleanup(cleanup)

	fx := fixtures.New(database)
	require.NoError(t, fx.Reset(ctx))
	require.NoError(t, fx.SeedAll(ctx))

	ldap := &mockLdap{}

	uc := NewUseCase(
		NewFactory(database, zap.NewNop()),
		ldap,
		zap.NewNop(),
	)

	return &testEnv{
		uc:   uc,
		ctx:  ctx,
		repo: userrepo.NewRepository(database, zap.NewNop()),
	}
}

// =========================
// TESTS
// =========================

func TestCreateLecturer_Success(t *testing.T) {
	env := setup(t)

	ldap := &mockLdap{
		GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
			return &ldapapdap.LdapUserInfo{
				Login: login,
				Name:  "John Doe",
				Mail:  "john@test.com",
			}, nil
		},
	}

	env.uc = NewUseCase(env.uc.factory, ldap, zap.NewNop())

	req := CreateLecturerRequest{
		AdminID: 8, // olgbvl
		Login:   "new_lecturer",
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestCreateLecturer_InvalidInput(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, CreateLecturerRequest{
		AdminID: 0,
		Login:   "",
	})

	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestCreateLecturer_Forbidden_NotAdmin(t *testing.T) {
	env := setup(t)

	ldap := &mockLdap{
		GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
			return &ldapapdap.LdapUserInfo{
				Login: login,
				Name:  "John Doe",
				Mail:  "john@test.com",
			}, nil
		},
	}

	env.uc = NewUseCase(env.uc.factory, ldap, zap.NewNop())

	req := CreateLecturerRequest{
		AdminID: 1, // student → не admin
		Login:   "new_lecturer",
	}

	_, err := env.uc.Execute(env.ctx, req)

	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateLecturer_LDAPError(t *testing.T) {
	env := setup(t)

	ldap := &mockLdap{
		GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
			return nil, errors.New("ldap down")
		},
	}

	env.uc = NewUseCase(env.uc.factory, ldap, zap.NewNop())

	req := CreateLecturerRequest{
		AdminID: 8,
		Login:   "new_lecturer",
	}

	_, err := env.uc.Execute(env.ctx, req)

	assert.ErrorIs(t, err, ErrLDAPFailed)
}

func TestCreateLecturer_Idempotent(t *testing.T) {
	env := setup(t)

	ldap := &mockLdap{
		GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
			return &ldapapdap.LdapUserInfo{
				Login: "lecturer1", // уже существует
				Name:  "John Doe",
				Mail:  "john@test.com",
				Group: nil,
			}, nil
		},
	}

	env.uc = NewUseCase(env.uc.factory, ldap, zap.NewNop())

	req := CreateLecturerRequest{
		AdminID: 8,
		Login:   "lecturer1",
	}

	resp, err := env.uc.Execute(env.ctx, req)

	fmt.Println(resp, err)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

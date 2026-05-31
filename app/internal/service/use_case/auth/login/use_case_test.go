package login

import (
	"context"
	"errors"
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
	ValidateFunc func(ctx context.Context, login, password string) error
	GetInfoFunc  func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error)
}

func (m *mockLDAP) ValidatePassword(ctx context.Context, login, password string) error {
	return m.ValidateFunc(ctx, login, password)
}

func (m *mockLDAP) GetInfo(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
	return m.GetInfoFunc(ctx, login)
}

//
// =========================
// MOCK AUTH
// =========================
//

type mockAuth struct {
	GenAccessFunc  func(userID int) (string, error)
	GenRefreshFunc func(userID int) (string, error)
}

func (m *mockAuth) GenerateAccess(userID int) (string, error) {
	return m.GenAccessFunc(userID)
}

func (m *mockAuth) GenerateRefresh(userID int) (string, error) {
	return m.GenRefreshFunc(userID)
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
			nil,
			nil,
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

func TestAuthLogin_Success_WithFixtures(t *testing.T) {
	env := setup(t)

	ldap := &mockLDAP{
		ValidateFunc: func(ctx context.Context, login, password string) error {
			return nil
		},
		GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
			return &ldapapdap.LdapUserInfo{
				Login: "student1",
				Name:  "Student One",
				Mail:  "student1@mail.com",
			}, nil
		},
	}

	auth := &mockAuth{
		GenAccessFunc: func(userID int) (string, error) {
			return "access-token", nil
		},
		GenRefreshFunc: func(userID int) (string, error) {
			return "refresh-token", nil
		},
	}

	env.uc = NewUseCase(
		NewFactory(env.uc.factory.(*factory).db, zap.NewNop()),
		ldap,
		auth,
		zap.NewNop(),
	)

	resp, err := env.uc.Execute(env.ctx, AuthLoginRequest{
		Login:    "student1",
		Password: "123",
	})

	require.NoError(t, err)
	assert.Equal(t, "access-token", resp.AccessToken)
	assert.Equal(t, "refresh-token", resp.RefreshToken)
}

func TestAuthLogin_InvalidInput(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, AuthLoginRequest{
		Login: "",
	})

	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestAuthLogin_LDAPValidationFail(t *testing.T) {
	env := setup(t)

	env.uc = NewUseCase(
		env.uc.factory,
		&mockLDAP{
			ValidateFunc: func(ctx context.Context, login, password string) error {
				return errors.New("invalid credentials")
			},
		},
		&mockAuth{},
		zap.NewNop(),
	)

	_, err := env.uc.Execute(env.ctx, AuthLoginRequest{
		Login:    "student1",
		Password: "wrong",
	})

	assert.ErrorIs(t, err, ErrUnauthorized)
}

func TestAuthLogin_LDAPGetInfoFail(t *testing.T) {
	env := setup(t)

	env.uc = NewUseCase(
		env.uc.factory,
		&mockLDAP{
			ValidateFunc: func(ctx context.Context, login, password string) error {
				return nil
			},
			GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
				return nil, errors.New("ldap down")
			},
		},
		&mockAuth{},
		zap.NewNop(),
	)

	_, err := env.uc.Execute(env.ctx, AuthLoginRequest{
		Login:    "student1",
		Password: "123",
	})

	assert.ErrorIs(t, err, ErrLDAPFailed)
}

func TestAuthLogin_TokenGenerationFail(t *testing.T) {
	env := setup(t)

	env.uc = NewUseCase(
		env.uc.factory,
		&mockLDAP{
			ValidateFunc: func(ctx context.Context, login, password string) error {
				return nil
			},
			GetInfoFunc: func(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error) {
				return &ldapapdap.LdapUserInfo{
					Login: "student1",
					Name:  "Student One",
					Mail:  "student1@mail.com",
				}, nil
			},
		},
		&mockAuth{
			GenAccessFunc: func(userID int) (string, error) {
				return "", errors.New("fail")
			},
			GenRefreshFunc: func(userID int) (string, error) {
				return "refresh", nil
			},
		},
		zap.NewNop(),
	)

	_, err := env.uc.Execute(env.ctx, AuthLoginRequest{
		Login:    "student1",
		Password: "123",
	})

	assert.ErrorIs(t, err, ErrTokenGenerate)
}

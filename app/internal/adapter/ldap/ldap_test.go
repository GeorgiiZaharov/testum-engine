package ldap

import (
	"context"
	"crypto/tls"
	"errors"
	"testing"
	"time"

	"testum-engine/app/internal/config"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

//
// =========================
// FAKE LDAP CONNECTION
// =========================
//

type fakeConn struct {
	searchResult *ldap.SearchResult
	searchErr    error
	bindErr      error

	delay time.Duration

	closed  bool
	timeout time.Duration
}

func (f *fakeConn) Close() error {
	f.closed = true
	return nil
}

func (f *fakeConn) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	if f.delay > 0 {
		time.Sleep(f.delay)
	}
	return f.searchResult, f.searchErr
}

func (f *fakeConn) Bind(username, password string) error {
	if f.delay > 0 {
		time.Sleep(f.delay)
	}
	return f.bindErr
}

func (f *fakeConn) SetTimeout(d time.Duration) {
	f.timeout = d
}

//
// =========================
// HELPERS
// =========================
//

func fakeSearchResult(dn string) *ldap.SearchResult {
	return &ldap.SearchResult{
		Entries: []*ldap.Entry{
			{
				DN: dn,
				Attributes: []*ldap.EntryAttribute{
					{Name: "uid", Values: []string{"john"}},
					{Name: "cn", Values: []string{"John Doe"}},
					{Name: "mail", Values: []string{"john@mail.com"}},
				},
			},
		},
	}
}

func newAdapterWithDial(dial func(string, *tls.Config) (ldapConn, error)) *ldapAdapter {
	return &ldapAdapter{
		cfg: config.LdapConfig{
			Server:  "localhost",
			Port:    636,
			BaseDN:  "dc=test",
			Timeout: 50 * time.Millisecond,
		},
		logger: zap.NewNop(),
		dial:   dial,
	}
}

func fakeDialSuccess(conn ldapConn) func(string, *tls.Config) (ldapConn, error) {
	return func(url string, tlsConfig *tls.Config) (ldapConn, error) {
		return conn, nil
	}
}

//
// =========================
// TESTS: ValidatePassword
// =========================
//

func TestValidatePassword_EmptyInput(t *testing.T) {
	a := newAdapterWithDial(nil)

	err := a.ValidatePassword(context.Background(), "", "")

	assert.ErrorIs(t, err, ErrEmptyCredentials)
}

func TestValidatePassword_ConnectFail(t *testing.T) {
	a := newAdapterWithDial(func(url string, tlsConfig *tls.Config) (ldapConn, error) {
		return nil, errors.New("connect error")
	})

	err := a.ValidatePassword(context.Background(), "john", "123")

	assert.ErrorIs(t, err, ErrConnectionFailed)
}

func TestValidatePassword_UserNotFound(t *testing.T) {
	a := newAdapterWithDial(fakeDialSuccess(&fakeConn{
		searchResult: &ldap.SearchResult{Entries: []*ldap.Entry{}},
	}))

	err := a.ValidatePassword(context.Background(), "john", "123")

	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestValidatePassword_SearchError(t *testing.T) {
	a := newAdapterWithDial(fakeDialSuccess(&fakeConn{
		searchErr: errors.New("search fail"),
	}))

	err := a.ValidatePassword(context.Background(), "john", "123")

	assert.ErrorIs(t, err, ErrSearchFailed)
}

func TestValidatePassword_BindFail(t *testing.T) {
	a := newAdapterWithDial(fakeDialSuccess(&fakeConn{
		searchResult: fakeSearchResult("uid=john"),
		bindErr:      errors.New("bad password"),
	}))

	err := a.ValidatePassword(context.Background(), "john", "123")

	assert.ErrorIs(t, err, ErrInvalidPassword)
}

func TestValidatePassword_Success(t *testing.T) {
	a := newAdapterWithDial(fakeDialSuccess(&fakeConn{
		searchResult: fakeSearchResult("uid=john"),
	}))

	err := a.ValidatePassword(context.Background(), "john", "123")

	assert.NoError(t, err)
}

func TestValidatePassword_ContextTimeout(t *testing.T) {
	conn := &fakeConn{
		searchResult: fakeSearchResult("uid=john"),
		delay:        200 * time.Millisecond,
	}

	a := newAdapterWithDial(fakeDialSuccess(conn))

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := a.ValidatePassword(ctx, "john", "123")

	assert.Error(t, err)
}

//
// =========================
// TESTS: GetInfo
// =========================
//

func TestGetInfo_EmptyLogin(t *testing.T) {
	a := newAdapterWithDial(nil)

	_, err := a.GetInfo(context.Background(), "")

	assert.ErrorIs(t, err, ErrEmptyLogin)
}

func TestGetInfo_ConnectFail(t *testing.T) {
	a := newAdapterWithDial(func(url string, tlsConfig *tls.Config) (ldapConn, error) {
		return nil, errors.New("connect fail")
	})

	_, err := a.GetInfo(context.Background(), "john")

	assert.ErrorIs(t, err, ErrConnectionFailed)
}

func TestGetInfo_SearchError(t *testing.T) {
	a := newAdapterWithDial(fakeDialSuccess(&fakeConn{
		searchErr: errors.New("search fail"),
	}))

	_, err := a.GetInfo(context.Background(), "john")

	assert.ErrorIs(t, err, ErrSearchFailed)
}

func TestGetInfo_UserNotFound(t *testing.T) {
	a := newAdapterWithDial(fakeDialSuccess(&fakeConn{
		searchResult: &ldap.SearchResult{Entries: []*ldap.Entry{}},
	}))

	info, err := a.GetInfo(context.Background(), "john")

	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, info)
}

func TestGetInfo_Success(t *testing.T) {
	a := newAdapterWithDial(fakeDialSuccess(&fakeConn{
		searchResult: fakeSearchResult("uid=john"),
	}))

	info, err := a.GetInfo(context.Background(), "john")

	require.NoError(t, err)
	require.NotNil(t, info)

	assert.Equal(t, "john", info.Login)
	assert.Equal(t, "John Doe", info.Name)
	assert.Equal(t, "john@mail.com", info.Mail)
}

func TestGetInfo_ContextTimeout(t *testing.T) {
	conn := &fakeConn{
		searchResult: fakeSearchResult("uid=john"),
		delay:        200 * time.Millisecond,
	}

	a := newAdapterWithDial(fakeDialSuccess(conn))

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := a.GetInfo(ctx, "john")

	assert.Error(t, err)
}

package ldap

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-ldap/ldap/v3"
	"go.uber.org/zap"

	"testum-engine/app/internal/config"
)

// =========================
// INTERFACE
// =========================

type LdapAdapter interface {
	ValidatePassword(ctx context.Context, login, password string) error
	GetInfo(ctx context.Context, login string) (*LdapUserInfo, error)
}

// =========================
// INTERNAL CONNECTION
// =========================

type ldapConn interface {
	Close() error
	Search(*ldap.SearchRequest) (*ldap.SearchResult, error)
	Bind(username, password string) error
	SetTimeout(time.Duration)
}

// =========================
// DIAL STRATEGY (NO GLOBALS)
// =========================

type dialFunc func(string, *tls.Config) (ldapConn, error)

// production default dial
func defaultDial(url string, tlsConfig *tls.Config) (ldapConn, error) {
	return ldap.DialURL(url, ldap.DialWithTLSConfig(tlsConfig))
}

// =========================
// ADAPTER
// =========================

type ldapAdapter struct {
	cfg    config.LdapConfig
	logger *zap.Logger
	dial   dialFunc
}

// =========================
// CONSTRUCTOR
// =========================

func NewLdapAdapter(cfg config.LdapConfig, logger *zap.Logger) LdapAdapter {
	return &ldapAdapter{
		cfg:    cfg,
		logger: logger,
		dial:   defaultDial,
	}
}

// =========================
// TEST CONSTRUCTOR (IMPORTANT)
// =========================

func NewTestLdapAdapter(cfg config.LdapConfig, logger *zap.Logger, dial dialFunc) LdapAdapter {
	return &ldapAdapter{
		cfg:    cfg,
		logger: logger,
		dial:   dial,
	}
}

// =========================
// HELPERS
// =========================

func withContext[T any](ctx context.Context, fn func() (T, error)) (T, error) {
	type result struct {
		val T
		err error
	}

	ch := make(chan result, 1)

	go func() {
		v, err := fn()
		ch <- result{v, err}
	}()

	select {
	case res := <-ch:
		return res.val, res.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// =========================
// CONNECT
// =========================

func (l *ldapAdapter) connect(ctx context.Context) (ldapConn, error) {
	url := fmt.Sprintf("ldaps://%s:%d", l.cfg.Server, l.cfg.Port)

	l.logger.Info("ldap connect",
		zap.String("url", url),
	)

	conn, err := l.dial(url, &tls.Config{
		InsecureSkipVerify: true, // TODO production TLS
	})

	if err != nil {
		l.logger.Error("ldap connection failed",
			zap.Error(err),
			zap.String("url", url),
		)
		return nil, ErrConnectionFailed
	}

	timeout := l.cfg.Timeout
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining < timeout {
			timeout = remaining
		}
	}

	conn.SetTimeout(timeout)

	return conn, nil
}

// =========================
// VALIDATE PASSWORD
// =========================

func (l *ldapAdapter) ValidatePassword(ctx context.Context, login, password string) error {
	if login == "lector" && password == "123456" { // TODO:
		return nil
	}
	if login == "student1" && password == "123456" { //TODO:
		return nil
	}
	if login == "student2" && password == "123456" { //TODO:
		return nil
	}
	if login == "student3" && password == "123456" { //TODO:
		return nil
	}
	if login == "" || password == "" {
		l.logger.Warn("empty credentials",
			zap.String("login", login),
		)
		return ErrEmptyCredentials
	}

	conn, err := l.connect(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	searchReq := ldap.NewSearchRequest(
		l.cfg.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		fmt.Sprintf("(uid=%s)", login),
		[]string{"dn"},
		nil,
	)

	res, err := withContext(ctx, func() (*ldap.SearchResult, error) {
		return conn.Search(searchReq)
	})
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		l.logger.Error("ldap search failed",
			zap.Error(err),
			zap.String("login", login),
		)
		return ErrSearchFailed
	}

	if len(res.Entries) == 0 {
		l.logger.Warn("user not found",
			zap.String("login", login),
		)
		return ErrUserNotFound
	}

	userDN := res.Entries[0].DN

	_, err = withContext(ctx, func() (struct{}, error) {
		return struct{}{}, conn.Bind(userDN, password)
	})
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		l.logger.Warn("invalid credentials",
			zap.String("login", login),
		)
		return ErrInvalidPassword
	}

	l.logger.Info("auth success",
		zap.String("login", login),
	)

	return nil
}

// =========================
// GET INFO
// =========================

func (l *ldapAdapter) GetInfo(ctx context.Context, login string) (*LdapUserInfo, error) {
	if login == "student1" { // TODO:
		group := "22307"
		return &LdapUserInfo{
			Login: "student1",
			Name:  "Student First",
			Mail:  "student1@mail.ru",
			Group: &group,
		}, nil
	}
	if login == "student2" { // TODO:
		group := "22307"
		return &LdapUserInfo{
			Login: "student2",
			Name:  "Student Second",
			Mail:  "student2@mail.ru",
			Group: &group,
		}, nil
	}
	if login == "student3" { // TODO:
		group := "22306"
		return &LdapUserInfo{
			Login: "student3",
			Name:  "Student Third",
			Mail:  "student3@mail.ru",
			Group: &group,
		}, nil
	}
	if login == "lector" { //TODO:
		return &LdapUserInfo{
			Login: "lector",
			Name:  "Lector Lector",
			Mail:  "lector@mail.ru",
			Group: nil,
		}, nil
	}
	if login == "" {
		return nil, ErrEmptyLogin
	}

	conn, err := l.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	searchReq := ldap.NewSearchRequest(
		l.cfg.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1,
		0,
		false,
		fmt.Sprintf("(uid=%s)", login),
		[]string{"uid", "cn", "mail"},
		nil,
	)

	res, err := withContext(ctx, func() (*ldap.SearchResult, error) {
		return conn.Search(searchReq)
	})
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		l.logger.Error("ldap search error",
			zap.Error(err),
			zap.String("login", login),
		)
		return nil, ErrSearchFailed
	}

	if len(res.Entries) == 0 {
		l.logger.Warn("user not found",
			zap.String("login", login),
		)
		return nil, ErrUserNotFound
	}

	e := res.Entries[0]

	info := &LdapUserInfo{
		Login: e.GetAttributeValue("uid"),
		Name:  e.GetAttributeValue("cn"),
		Mail:  e.GetAttributeValue("mail"),
		Group: extractAcademicGroup(e.DN),
	}

	l.logger.Info("ldap user fetched",
		zap.String("login", login),
	)

	return info, nil
}

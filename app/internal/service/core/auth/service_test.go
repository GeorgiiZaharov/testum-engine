package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func newService(t *testing.T) Service {
	logger := zaptest.NewLogger(t)
	return New("secret", logger)
}

func TestGenerateAccess(t *testing.T) {
	s := newService(t)

	tokenStr, err := s.GenerateAccess(42)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			return []byte("secret"), nil
		},
	)
	require.NoError(t, err)

	claims := token.Claims.(*Claims)

	assert.Equal(t, 42, claims.UserID)
	assert.Equal(t, TokenTypeAccess, claims.Type)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestGenerateRefresh(t *testing.T) {
	s := newService(t)

	tokenStr, err := s.GenerateRefresh(99)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			return []byte("secret"), nil
		},
	)
	require.NoError(t, err)

	claims := token.Claims.(*Claims)

	assert.Equal(t, 99, claims.UserID)
	assert.Equal(t, TokenTypeRefresh, claims.Type)
}

func TestParse_ValidToken(t *testing.T) {
	s := newService(t)

	tokenStr, err := s.GenerateAccess(1)
	require.NoError(t, err)

	claims, err := s.Parse(tokenStr)
	require.NoError(t, err)

	assert.Equal(t, 1, claims.UserID)
	assert.Equal(t, TokenTypeAccess, claims.Type)
}

func TestParse_InvalidTokenString(t *testing.T) {
	s := newService(t)

	_, err := s.Parse("not-a-real-token")
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestParse_WrongSecret(t *testing.T) {
	logger := zaptest.NewLogger(&testing.T{})

	s1 := New("secret1", logger)
	s2 := New("secret2", logger)

	token, err := s2.GenerateAccess(1)
	require.NoError(t, err)

	_, err = s1.Parse(token)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestParse_ExpiredToken(t *testing.T) {
	logger := zaptest.NewLogger(t)

	s := &service{
		secret: []byte("secret"),
		log:    logger,
	}

	claims := Claims{
		UserID: 10,
		Type:   TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	_, err = s.Parse(tokenStr)
	assert.ErrorIs(t, err, ErrExpiredToken)
}

func TestParse_InvalidClaimsCast(t *testing.T) {
	logger := zaptest.NewLogger(t)

	s := &service{
		secret: []byte("secret"),
		log:    logger,
	}

	// токен с НЕ тем типом claims
	type fakeClaims struct {
		jwt.RegisteredClaims
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, fakeClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})

	tokenStr, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	_, err = s.Parse(tokenStr)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestParse_InvalidTokenClaimsType(t *testing.T) {
	s := newService(t)

	// создаём токен, но ломаем подпись чтобы token.Valid = false
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})

	// подписываем НЕ тем секретом → token.Valid станет false
	tokenStr, err := token.SignedString([]byte("wrong-secret"))
	require.NoError(t, err)

	_, err = s.Parse(tokenStr)
	assert.ErrorIs(t, err, ErrInvalidToken)
}
func TestParse_InvalidClaimsContent(t *testing.T) {
	s := newService(t)

	claims := Claims{
		UserID: 0, // ❌ invalid
		Type:   "",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte("secret"))
	require.NoError(t, err)

	_, err = s.Parse(tokenStr)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestParse_InvalidTokenValidFalse(t *testing.T) {
	logger := zaptest.NewLogger(t)

	s1 := New("secret1", logger)
	s2 := New("secret2", logger)

	token, err := s2.GenerateAccess(123)
	require.NoError(t, err)

	_, err = s1.Parse(token)

	// именно этот кейс покрывает !token.Valid
	require.ErrorIs(t, err, ErrInvalidToken)
}

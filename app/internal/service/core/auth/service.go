package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Service interface {
	GenerateAccess(userID int) (string, error)
	GenerateRefresh(userID int) (string, error)
	Parse(tokenStr string) (*Claims, error)
}

type service struct {
	secret []byte
	log    *zap.Logger
}

func New(secret string, log *zap.Logger) Service {
	return &service{
		secret: []byte(secret),
		log:    log.Named("auth-jwt"),
	}
}

func (s *service) GenerateAccess(userID int) (string, error) {
	s.log.Debug("generating access token",
		zap.Int("user_id", userID),
	)

	claims := Claims{
		UserID: userID,
		Type:   TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.secret)
	if err != nil {
		s.log.Error("failed to sign access token",
			zap.Int("user_id", userID),
			zap.Error(err),
		)
		return "", err
	}

	s.log.Debug("access token generated",
		zap.Int("user_id", userID),
	)

	return signed, nil
}
func (s *service) GenerateRefresh(userID int) (string, error) {
	s.log.Debug("generating refresh token",
		zap.Int("user_id", userID),
	)

	claims := Claims{
		UserID: userID,
		Type:   TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.secret)
	if err != nil {
		s.log.Error("failed to sign refresh token",
			zap.Int("user_id", userID),
			zap.Error(err),
		)
		return "", err
	}

	s.log.Debug("refresh token generated",
		zap.Int("user_id", userID),
	)

	return signed, nil
}
func (s *service) Parse(tokenStr string) (*Claims, error) {
	s.log.Debug("parsing token")

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			return s.secret, nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			s.log.Info("expired token detected")
			return nil, ErrExpiredToken
		}

		s.log.Warn("invalid token parse error", zap.Error(err))
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		s.log.Warn("invalid token claims")
		return nil, ErrInvalidToken
	}
	if claims.UserID == 0 || claims.Type == "" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

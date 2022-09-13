package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ollub/user_service/internal/users"
	"github.com/Ollub/user_service/internal/users/usecase"
	"github.com/Ollub/user_service/pkg/log"
	"github.com/Ollub/user_service/pkg/utils"

	"github.com/dgrijalva/jwt-go"
)

var AuthError = errors.New("authentication error")

type SessionsJWTVer struct {
	Secret       []byte
	TokenTTLDays int
	users        *usecase.Manager
}

type SessionJWTVerClaims struct {
	UserID uint32 `json:"uid"`
	Ver    int    `json:"ver,omitempty"`
	jwt.StandardClaims
}

func NewSessionsJWTVer(secret []byte, tokenTTLDays int, manager *usecase.Manager) *SessionsJWTVer {
	return &SessionsJWTVer{
		Secret:       secret,
		TokenTTLDays: tokenTTLDays,
		users:        manager,
	}
}

func (sm *SessionsJWTVer) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return sm.Secret, nil
}

func (sm *SessionsJWTVer) Check(ctx context.Context, token string) (*Session, error) {
	payload := &SessionJWTVerClaims{}
	_, err := jwt.ParseWithClaims(token, payload, sm.parseSecretGetter)
	if err != nil {
		return nil, fmt.Errorf("cant parse jwt token: %v", err)
	}
	// check exp, iat
	if payload.Valid() != nil {
		return nil, fmt.Errorf("invalid jwt token: %v", err)
	}

	ver, err := sm.users.GetUserVersion(ctx, payload.UserID)
	if err != nil {
		log.Clog(ctx).Info("Authentication failed for user", log.Fields{"userId": payload.UserID, "error": err})
		return nil, AuthError
	}

	if payload.Ver != ver {
		log.Clog(ctx).Info(
			"Provided token with old user version",
			log.Fields{"userId": payload.UserID, "tokenVer": payload.Ver, "actualVer": ver},
		)
		return nil, AuthError
	}

	return &Session{
		ID:     payload.Id,
		UserID: payload.UserID,
	}, nil
}

func (sm *SessionsJWTVer) Create(ctx context.Context, user *users.User) (string, error) {
	data := SessionJWTVerClaims{
		UserID: user.ID,
		Ver:    user.Ver, // изменилось по сравнению со stateless-сессией
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(sm.TokenTTLDays) * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        utils.RandStringRunes(32),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString(sm.Secret)
}

package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mutoulbj/gocsms/internal/config"
	"github.com/mutoulbj/gocsms/internal/models"
	"github.com/mutoulbj/gocsms/internal/repository"
	"github.com/mutoulbj/gocsms/pkg/cache"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cache    *cache.Cache
	log      *logrus.Logger
	jwtCfg  *config.JWTConfig
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	TokenID   string `json:"token_id"`
	IsRefresh bool   `json:"is_refresh"`
	jwt.RegisteredClaims
}

func NewAuthService(
	userRepo *repository.UserRepository,
	cache *cache.Cache,
	log *logrus.Logger,
	jwtCfg *config.JWTConfig,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cache:    cache,
		log:      log,
		jwtCfg:  jwtCfg,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*TokenPair, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		s.log.Error("Failed to get user by username: ", err)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.log.Error("Invalid password for user: ", username)
		return nil, err
	}
	return s.generateTokenPair(user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshToken, true)
	if err != nil {
		s.log.Error("Failed to validate refresh token: ", err)
		return nil, err
	}
	// check if session is still valid in Redis
	sessionKey := "session:" + claims.UserID + ":" + claims.TokenID
	var sessionData string
	if err := s.cache.Get(ctx, sessionKey, &sessionData); err != nil {
		s.log.Warn("Session not found in Redis for user: ", claims.UserID)
		return nil, err
	}
	// parse user ID from claims
	userUUID, err := uuid.Parse(claims.UserID)
	if err != nil {
		s.log.Error("Failed to parse user ID to UUID: ", err)
		return nil, err
	}
	user, err := s.userRepo.GetUserById(ctx, userUUID)
	if err != nil {
		s.log.Error("Failed to get user by ID: ", err)
		return nil, err
	}

	return s.generateTokenPair(user)
}

func (s *AuthService) KickUser(ctx context.Context, userID, tokenID string) error {
	sessionKey := "session:" + userID + ":" + tokenID
	if err := s.cache.Delete(ctx, sessionKey); err != nil {
		return err
	}
	s.log.Info("User kicked successfully: ", userID)
	return nil
}

func (s *AuthService) Logout(ctx context.Context, userID, tokenID string) error {
	sessionKey := "session:" + userID + ":" + tokenID
	if err := s.cache.Delete(ctx, sessionKey); err != nil {
		return err
	}
	s.log.Info("User logged out successfully: ", userID)
	return nil
}

func (s *AuthService) generateTokenPair(user *models.User) (*TokenPair, error) {
	tokenID := generateTokenID()
	accessToken, err := s.generateJWT(user, tokenID, false)
	if err != nil {
		s.log.Error("Failed to generate access token: ", err)
		return nil, err
	}

	refreshToken, err := s.generateJWT(user, tokenID, true)
	if err != nil {
		s.log.Error("Failed to generate refresh token: ", err)
		return nil, err
	}

	// store session in cache
	sessionKey := "session:" + user.ID.String() + ":" + tokenID
	err = s.cache.Set(context.Background(), sessionKey, "active", s.jwtCfg.RefreshTokenTTL)
	if err != nil {
		s.log.Error("Failed to store session in cache: ", err)
		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateJWT(user *models.User, tokenID string, isRefresh bool) (string, error) {
	expiresIn := s.jwtCfg.AccessTokenTTL
	if isRefresh {
		expiresIn = s.jwtCfg.RefreshTokenTTL
	}
	claims := Claims{
		UserID:    user.ID.String(),
		Username:  user.Username,
		TokenID:   tokenID,
		IsRefresh: isRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
			Issuer:    s.jwtCfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtCfg.Secret))
}

func (s *AuthService) ValidateToken(tokenString string, isRefresh bool) (*Claims, error) {
	s.log.Debug("Validating token: ", tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtCfg.Secret), nil
	})

	if err != nil {
		s.log.Error("Failed to parse token: ", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.IsRefresh != isRefresh {
			return nil, jwt.ErrTokenInvalidClaims
		}
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

func generateTokenID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

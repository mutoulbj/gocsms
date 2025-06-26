package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mutoulbj/gocsms/internal/config"
)

type CsmsClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

var appConfig *config.Config

func SetAuthConfig(cfg *config.Config) {
	appConfig = cfg
}

func GenerateTokens(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}

	// Access Token
	td.AtExpires = time.Now().Add(appConfig.JWTAccessTokenTTL).Unix()
	td.AccessUUID = fmt.Sprintf("%s-%d", userID, time.Now().UnixNano())

	atClaims := CsmsClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.AtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        td.AccessUUID,
			Subject:   userID,
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	td.AccessToken, err = at.SignedString([]byte(appConfig.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh Token
	td.RtExpires = time.Now().Add(appConfig.JWTRefreshTokenTTL).Unix()
	td.RefreshUUID = fmt.Sprintf("%s-%d", userID, time.Now().UnixNano())
	rtClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(td.RtExpires, 0)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        td.RefreshUUID,
		Subject:   userID,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(appConfig.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return td, nil
}

func VerifyToken(tokenString string) (*CsmsClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CsmsClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(appConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*CsmsClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

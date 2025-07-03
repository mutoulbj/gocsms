package jwt

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

func GenerateTokens(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(appConfig.JWT.AccessTokenTTL).Unix()

	// Access Token
	td.AtExpires = time.Now().Add(appConfig.JWT.AccessTokenTTL).Unix()
	td.AccessUUID = fmt.Sprintf("%s-%d", userID, time.Now().UnixNano())

	atClaims := CsmsClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.AtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        td.AccessUUID,
			Subject:   userID,
			Issuer:    appConfig.JWT.Issuer,
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	td.AccessToken, err = at.SignedString([]byte(appConfig.JWT.Secret))
	if err != nil {
		td.RtExpires = time.Now().Add(appConfig.JWT.RefreshTokenTTL).Unix()
	}

	// Refresh Token
	td.RtExpires = time.Now().Add(appConfig.JWT.RefreshTokenTTL).Unix()
	td.RefreshUUID = fmt.Sprintf("%s-%d", userID, time.Now().UnixNano())
	rtClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(td.RtExpires, 0)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        td.RefreshUUID,
		Subject:   userID,
		Issuer:    appConfig.JWT.Issuer,
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(appConfig.JWT.Secret))
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
		return []byte(appConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*CsmsClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

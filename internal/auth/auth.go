package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

// HashPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// CheckPasswordHash -
func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}
	parts := strings.Fields(authHeader)

	if len(parts) != 2 {
		return "", errors.New("invalid authorization token")
	}

	if !strings.EqualFold(parts[0], "bearer") || len(parts[1]) == 0 {
		return "", errors.New("invalid authorization token")
	}
	return parts[1], nil
}

func GetApiKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", errors.New("ApiKey is missing")
	}
	parts := strings.Fields(apiKey)
	return parts[1], nil
}

func MakeRefreshToken() (string, error) {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return "", errors.New("can't generate refresh token")
	}
	return hex.EncodeToString(data), nil
}

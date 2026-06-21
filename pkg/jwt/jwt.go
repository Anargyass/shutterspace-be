// Package jwt menyediakan fungsi sign dan verify JWT token.
// SECURITY: Algoritma dikunci ke HS256. Token 'none' algorithm ditolak.
// SECURITY: Secret key wajib dari environment variable, tidak ada hardcoded fallback.
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims adalah struktur klaim JWT custom
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// SignToken membuat JWT token baru dengan klaim yang diberikan.
// expiryHours menentukan berapa jam token berlaku.
func SignToken(userID uuid.UUID, email, role, secret string, expiryHours int) (string, error) {
	claims := &Claims{
		UserID: userID.String(),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "shutterspace-api",
		},
	}

	// Hardcode HS256 — tidak mengikuti algoritma dari header token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// VerifyToken memvalidasi dan mem-parse JWT token.
// Mengembalikan Claims jika valid, error jika tidak.
func VerifyToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// SECURITY: Tolak semua algoritma selain HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}), // Whitelist hanya HS256
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

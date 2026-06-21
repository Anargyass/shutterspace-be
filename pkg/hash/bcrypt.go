// Package hash menyediakan fungsi hashing password menggunakan bcrypt.
// TODO(security): Pertimbangkan upgrade ke Argon2id untuk produksi.
package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// costFactor adalah bcrypt cost factor. Nilai 12 sudah cukup aman untuk prototipe.
const costFactor = 12

// HashPassword meng-hash password menggunakan bcrypt dengan cost factor 12.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), costFactor)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword membandingkan password plaintext dengan hash yang tersimpan.
// Mengembalikan nil jika cocok, error jika tidak.
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

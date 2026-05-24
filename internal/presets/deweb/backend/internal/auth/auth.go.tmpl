package auth

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var loginRe = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func ValidateLogin(login string) error {
	login = strings.TrimSpace(login)
	if utf8.RuneCountInString(login) < 6 {
		return errors.New("логин: минимум 6 символов, только латиница и цифры")
	}
	if !loginRe.MatchString(login) {
		return errors.New("логин: только латинские буквы и цифры")
	}
	return nil
}

func ValidatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return errors.New("пароль: минимум 8 символов")
	}
	return nil
}

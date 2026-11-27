package service

import (
	"regexp"
	"strings"

	"github.com/ananda1farelas/sewakos-backend/helper"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateRegisterInput(name, email, password string) *helper.AppError {
	errs := map[string]string{}

	name = strings.TrimSpace(name)
	if len(name) < 3 {
		errs["name"] = "nama minimal 3 karakter"
	}

	if !emailRegex.MatchString(email) {
		errs["email"] = "format email tidak valid"
	}

	if len(password) < 8 {
		errs["password"] = "password minimal 8 karakter"
	}

	if len(errs) > 0 {
		return helper.ErrValidation(errs)
	}

	return nil
}

func ValidateLoginInput(email, password string) *helper.AppError {
	errs := map[string]string{}

	if !emailRegex.MatchString(email) {
		errs["email"] = "format email tidak valid"
	}

	if password == "" {
		errs["password"] = "password wajib diisi"
	}

	if len(errs) > 0 {
		return helper.ErrValidation(errs)
	}

	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

const DefaultRole = "tenant"

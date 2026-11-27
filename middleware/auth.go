package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ananda1farelas/sewakos-backend/config"
	"github.com/ananda1farelas/sewakos-backend/helper"
)

func RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return helper.ErrUnauthorized("token tidak ditemukan")
	}

	tokenString := authHeader[7:]

	claims, err := helper.ParseAccessToken(tokenString, config.Env.JWTSecret)
	if err != nil {
		return helper.ErrUnauthorized("token tidak valid atau kedaluwarsa")
	}

	helper.SetAuthLocals(c, claims.UserID, claims.Role)

	return c.Next()
}

// rate limiter login: 5 kali per menit per IP, balikin 429
type loginAttempt struct {
	count     int
	expiresAt time.Time
}

var (
	loginAttempts = make(map[string]*loginAttempt)
	loginMu       sync.Mutex
)

func LoginRateLimiter(c *fiber.Ctx) error {
	ip := c.IP()

	loginMu.Lock()
	defer loginMu.Unlock()

	now := time.Now()
	attempt, exists := loginAttempts[ip]

	if !exists || now.After(attempt.expiresAt) {
		loginAttempts[ip] = &loginAttempt{count: 1, expiresAt: now.Add(1 * time.Minute)}
		return c.Next()
	}

	if attempt.count >= 5 {
		return helper.NewAppError(helper.CodeForbidden, "terlalu banyak percobaan login, coba lagi nanti", 429, nil)
	}

	attempt.count++
	return c.Next()
}

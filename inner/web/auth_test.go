package web

import (
	jwtMiddleware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"
)

const testSecret = "super-secret-for-tests"

func buildApp() *fiber.App {
	app := fiber.New()

	// JWT middleware
	app.Use(jwtMiddleware.New(jwtMiddleware.Config{
		SigningKey: jwtMiddleware.SigningKey{Key: []byte(testSecret)},
		// Разрешаем запросы без токена
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Любая JWT-ошибка → 401
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	}))

	// Ролевая проверка
	app.Use(func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		rolesRaw, ok := claims["realm_access"]
		if !ok {
			return c.SendStatus(fiber.StatusForbidden)
		}

		realmMap, ok := rolesRaw.(map[string]interface{})
		if !ok {
			return c.SendStatus(fiber.StatusForbidden)
		}
		rolesRaw, ok = realmMap["roles"]
		if !ok {
			return c.Status(fiber.StatusForbidden).SendString("no roles")
		}
		// приводим []interface{} -> []string
		rolesSl, ok := rolesRaw.([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).SendString("roles is not an array")
		}
		var roles []string
		for _, r := range rolesSl {
			if s, ok := r.(string); ok {
				roles = append(roles, s)
			}
		}

		if slices.Contains(roles, IdmAdmin) {
			return c.Next()
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{})
	})

	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendString("admin area")
	})

	return app
}

func TestJWT_Auth(t *testing.T) {
	var a = assert.New(t)
	server := buildApp()

	test := []struct {
		name   string
		token  string
		status int
	}{
		// 401
		{"no_token", "", http.StatusUnauthorized},
		{"malformed", "Bearer INVALID_TOKEN", http.StatusUnauthorized},
		{"expired", buildToken(-time.Hour, []string{IdmAdmin}), http.StatusUnauthorized},
		{"wrong_sign", buildToken(time.Hour, []string{IdmAdmin}, "bad-secret"), http.StatusUnauthorized},

		// 403
		{"no_roles", buildToken(time.Hour, []string{}), http.StatusForbidden},
		{"wrong_role", buildToken(time.Hour, []string{"user"}), http.StatusForbidden},
		{"wrong_role", buildToken(time.Hour, []string{IdmAdmin}), http.StatusOK},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			resp, err := server.Test(req)
			a.NoError(err)
			a.Equal(tt.status, resp.StatusCode)
		})
	}

}

func buildToken(exp time.Duration, roles []string, secret ...string) string {
	key := testSecret
	if len(secret) > 0 {
		key = secret[0]
	}

	idmClaims := IdmClaims{
		RealmAccess: RealmAccessClaims{
			Roles: roles,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, idmClaims).SignedString([]byte(key))
	return "Bearer " + token
}

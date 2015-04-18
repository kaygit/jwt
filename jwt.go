package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrMissingConfig   = errors.New("missing configuration")
	ErrMissingSecret   = errors.New("please provide a shared secret")
	ErrMissingAuthFunc = errors.New("please provide an auth function")
)

type Config struct {
	Secret string
	Auth   AuthFunc
}

type AuthFunc func(string, string) (bool, error)

type JWTMiddleware struct {
	secret string
	auth   AuthFunc
}

func NewMiddleware(c *Config) (*JWTMiddleware, error) {
	if c == nil {
		return nil, ErrMissingConfig
	}
	if c.Secret == "" {
		return nil, ErrMissingSecret
	}
	if c.Auth == nil {
		return nil, ErrMissingAuthFunc
	}
	m := &JWTMiddleware{
		secret: c.Secret,
		auth:   c.Auth,
	}
	return m, nil
}

func (m *JWTMiddleware) Secure(h http.Handler) http.Handler {
	// This is just a placeholder for now
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (m *JWTMiddleware) GenerateToken(w http.ResponseWriter, r *http.Request) {
	var b map[string]string
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		panic(err)
	}
	result, err := m.auth(b["email"], b["password"])
	if err != nil {
		panic(err)
	}
	if !result {
		panic("deal with this")
	}

	// For now, the header will be static
	resp := `{"typ":"JWT","alg":"HS256"}`
	resp = base64.StdEncoding.EncodeToString([]byte(resp))

	w.Write([]byte(resp))
}

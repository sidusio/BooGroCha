package srv

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"sidus.io/boogrocha/internal/credentials"
)

const (
	credentialsCookieKey        = "credentials"
	credentialsContextKey       = "credentials"
	credentialsCookieExpiryTime = time.Hour * 24 * 30
)

func (s *server) saveCredentials(w http.ResponseWriter, r *http.Request) {
	jReq := struct {
		CID      string
		Password string
	}{}

	err := json.NewDecoder(r.Body).Decode(&jReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cipherText, err := credentials.Encrypt(
		credentials.Credentials{
			CID:      jReq.CID,
			Password: jReq.Password,
		},
		s.credentialsSecret,
	)

	http.SetCookie(w, &http.Cookie{
		Name:     credentialsCookieKey,
		Value:    base64.StdEncoding.EncodeToString(cipherText),
		Path:     "/api",
		Domain:   "",
		Expires:  time.Now().Add(credentialsCookieExpiryTime),
		Secure:   false, // TODO: Should be true in production
		HttpOnly: true,
	})
}

func (*server) clearCredentials(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     credentialsCookieKey,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	})
}

func (s *server) middlewareExtractCredentials(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(credentialsCookieKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ciphertext, err := base64.StdEncoding.DecodeString(c.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		creds, err := credentials.Decrypt(ciphertext, s.credentialsSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), credentialsContextKey, creds)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) middlewareRefreshCredentials(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(credentialsCookieKey)
		if err == nil {
			http.SetCookie(w, &http.Cookie{
				Name:     credentialsCookieKey,
				Value:    c.Value,
				Path:     "/api",
				Domain:   "",
				Expires:  time.Now().Add(credentialsCookieExpiryTime),
				Secure:   false, // TODO: Should be true in production
				HttpOnly: true,
			})
		}
		next.ServeHTTP(w, r)
	})
}

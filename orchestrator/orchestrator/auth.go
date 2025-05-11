package orchestrator

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/ZolotarevAlexandr/yl_final/db"
)

type ctxKey string

var (
	keyUserID ctxKey = "userID"
	jwtKey           = []byte("your-very-secret-key")
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ Login, Password string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := db.User{Login: req.Login, Password: string(hash)}
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "login already taken", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ Login, Password string }
	json.NewDecoder(r.Body).Decode(&req)
	var user db.User
	if err := db.DB.Where("login = ?", req.Login).First(&user).Error; err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	// issue JWT
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(jwtKey)
	json.NewEncoder(w).Encode(map[string]string{"token": ss})
}

// JWTAuthMiddleware protects HTTP endpoints, sets userID in context
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}
			return []byte(jwtKey), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), keyUserID, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext extracts the authenticated user ID from the context.
// Returns 0 if not present.
func UserIDFromContext(ctx context.Context) uint {
	if v := ctx.Value(keyUserID); v != nil {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

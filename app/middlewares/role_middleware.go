package middlewares

import (
	"net/http"
	"slices"

	"github.com/gieart87/gotoko/app/core/session/auth"
	"gorm.io/gorm"
)

func RoleMiddleware(next http.HandlerFunc, db *gorm.DB, roles ...string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := auth.CurrentUser(db, w, r)
		if !slices.Contains(roles, user.Role.Name) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

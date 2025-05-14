// middleware/auth.go
package middleware

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

// BasicAuthWithDB devuelve un middleware de autenticación básica usando la conexión db.
func BasicAuthWithDB(db *sql.DB) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Buscar usuario por nombre (o email) en BD
		var hashed string
		err := db.QueryRow("SELECT password FROM users WHERE email = ?", username).Scan(&hashed)
		if err != nil {
			return false, nil // usuario no existe
		}
		// Comparar contraseñas
		if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
			return false, nil
		}
		return true, nil
	})

}

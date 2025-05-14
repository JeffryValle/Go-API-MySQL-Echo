// handlers/user.go
package handlers

import (
	"github/JeffryValle/db"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // omitimos la contraseña en las respuestas JSON
}

// RegisterUser crea un nuevo usuario.
// Espera JSON { "name": "...", "email": "...", "password": "..." }.
func RegisterUser(c echo.Context) error {

	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Solicitud inválida"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error al cifrar la contraseña"})
	}

	res, err := db.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", u.Name, u.Email, hash)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error al crear el usuario"})
	}

	id, _ := res.LastInsertId()
	u.ID = id
	u.Password = ""
	return c.JSON(http.StatusCreated, u)
}

// LoginUser autentica un usuario existente.
// Espera JSON { "email": "...", "password": "..." }.
func LoginUser(c echo.Context) error {
	usuario := new(User)
	if err := c.Bind(usuario); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Solicitud inválida"})
	}
	// Buscar usuario por email
	row := db.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", usuario.Email)
	var u User
	var hashedPassword string
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &hashedPassword); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Credenciales inválidas"})
	}
	// Comparar contraseñas con1. bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(usuario.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Credenciales inválidas"})
	}
	u.Password = "" // no devolvemos la contraseña
	return c.JSON(http.StatusOK, u)
}

// GetUsers devuelve la lista de todos los usuarios (sin contraseñas).
func GetUsers(c echo.Context) error {
	rows, err := db.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error al obtener usuarios"})
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error al procesar usuarios"})
		}
		users = append(users, u)
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser devuelve un usuario por ID (parámetro de ruta).
func GetUser(c echo.Context) error {
	id := c.Param("id")
	row := db.DB.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)
	var u User
	if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Usuario no encontrado"})
	}
	return c.JSON(http.StatusOK, u)
}

// UpdateUser modifica los datos de un usuario existente.
func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Solicitud inválida"})
	}
	// Cifrar nueva contraseña si se proporciona
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error al cifrar la contraseña"})
	}
	// Ejecutar actualización SQL
	_, err = db.DB.Exec("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?", u.Name, u.Email, hashed, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error al actualizar el usuario"})
	}
	u.ID, _ = strconv.ParseInt(id, 10, 64)
	u.Password = ""
	return c.JSON(http.StatusOK, u)
}

// DeleteUser elimina un usuario por ID.
func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	_, err := db.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error al eliminar el usuario"})
	}
	return c.NoContent(http.StatusNoContent)
}

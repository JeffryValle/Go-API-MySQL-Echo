// routes/routes.go
package routes

import (
	"github/JeffryValle/db"
	"github/JeffryValle/handlers"
	"github/JeffryValle/middleware"

	"github.com/labstack/echo/v4"
)

// SetupRoutes configura todas las rutas de la API.
func SetupRoutes(e *echo.Echo) {
	// Rutas sin autenticación
	e.POST("/register", handlers.RegisterUser)
	e.POST("/login", handlers.LoginUser)

	// Grupo protegido: requiere Basic Auth para /users...
	g := e.Group("/users")

	g.Use(middleware.BasicAuthWithDB(db.DB))

	g.GET("", handlers.GetUsers)
	g.GET("/:id", handlers.GetUser)
	g.PUT("/:id", handlers.UpdateUser)
	g.DELETE("/:id", handlers.DeleteUser)
}

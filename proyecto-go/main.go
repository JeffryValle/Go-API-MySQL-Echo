package main

import (
	"github/JeffryValle/db"
	"github/JeffryValle/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// Inicializar conexión a la base de datos
	db.Init()
	defer db.CloseConnection()

	// Crear instancia de Echo
	e := echo.New()

	// Registrar rutas
	routes.SetupRoutes(e)

	// Iniciar servidor en el puerto 1323
	e.Logger.Fatal(e.Start(":8080"))
}

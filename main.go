package main

import (
	"log"
	"shipping-app/internal/app/infrastructure/database/config"
	"shipping-app/internal/app/infrastructure/fiber/routers"
	"shipping-app/internal/utils"
)

func main() {

	db, err := config.SetDB()
	if err != nil {
		log.Fatalf("error al configurar la BD: %v", err)
		return
	}
	defer db.Close()

	app, err := utils.InitFiber()
	if err != nil {
		log.Fatal("error al configurar fiber")
	}

	routers.SetupRouters(app, db)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("error al iniciar el servidor: %v", err)
	}
}

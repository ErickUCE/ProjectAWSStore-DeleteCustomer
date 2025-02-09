package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ProjectAWSStore-DeleteCustomer/config"
	"ProjectAWSStore-DeleteCustomer/routes"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🚀 Iniciando DeleteCustomerService en Golang...")

	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ No se pudo cargar el archivo .env")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("❌ Error al conectar con MongoDB:", err)
	}
	fmt.Println("✅ Conexión exitosa a MongoDB")

	router := routes.SetupRoutes(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	fmt.Println("✅ Servidor corriendo en el puerto", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

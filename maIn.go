package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ProjectAWSStore-DeleteCustomer/config"
	"ProjectAWSStore-DeleteCustomer/routes"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🚀 Iniciando DeleteCustomerService en Golang...")

	// 📌 Cargar variables de entorno desde `.env`
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ Advertencia: No se pudo cargar el archivo .env, verificando variables de entorno...")
	}

	// 📌 Mostrar variables de entorno cargadas
	fmt.Println("🔗 MONGO_URI desde .env:", os.Getenv("MONGO_URI"))
	fmt.Println("🔗 MONGO_DB_NAME desde .env:", os.Getenv("MONGO_DB_NAME"))

	// 🔥 Si las variables están vacías, asignarlas manualmente
	if os.Getenv("MONGO_URI") == "" {
		fmt.Println("⚠️ No se encontró MONGO_URI, asignando manualmente...")
		os.Setenv("MONGO_URI", "mongodb://54.158.252.115:27017/CustomerDB")
		os.Setenv("MONGO_DB_NAME", "CustomerDB")
		os.Setenv("PORT", "8084")
	}

	// 📌 Verificar nuevamente después de forzar la carga
	fmt.Println("✅ MONGO_URI en uso:", os.Getenv("MONGO_URI"))

	// ✅ Conectar a MongoDB
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("❌ Error al conectar con MongoDB:", err)
	}
	fmt.Println("✅ Conexión exitosa a MongoDB")

	// ✅ Configurar rutas después de conectar a MongoDB
	router := routes.SetupRoutes(db)

	// 📌 Middleware CORS para permitir conexiones desde `http://localhost:3000`
	handler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Permitir solo el frontend
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(router)

	// 📌 Obtener el puerto desde `.env`
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // 🔥 Valor por defecto
	}

	// ✅ Iniciar el servidor en el puerto definido en `.env`
	fmt.Println("✅ Servidor owose corriendo en el puerto", port)
	log.Fatal(http.ListenAndServe(":"+port, handler)) // ✅ Ahora sí usa CORS
}

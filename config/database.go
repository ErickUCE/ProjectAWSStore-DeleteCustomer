package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Variable global para la conexión a la base de datos
var DB *mongo.Database

// ConnectDB establece la conexión con MongoDB
func ConnectDB() (*mongo.Database, error) {
	fmt.Println("📌 Ejecutando ConnectDB...")

	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ No se pudo cargar el archivo .env, verificando variables de entorno...")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("❌ Error: No se encontró MONGO_URI en las variables de entorno")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("❌ Error al conectar con MongoDB:", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Error: No se pudo conectar a MongoDB:", err)
		return nil, err
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "DeleteCustomerDB"
	}

	fmt.Println("🛢️ Base de datos seleccionada:", dbName)

	DB = client.Database(dbName)
	return DB, nil
}

// GetDB devuelve la instancia de la base de datos ya conectada
func GetDB() *mongo.Database {
	if DB == nil {
		log.Fatal("❌ Error: La base de datos no está inicializada. Asegúrate de llamar a ConnectDB() primero.")
	}
	return DB
}

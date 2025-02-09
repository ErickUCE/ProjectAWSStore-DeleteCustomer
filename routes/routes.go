package routes

import (
	"ProjectAWSStore-DeleteCustomer/controllers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes configura las rutas del API
func SetupRoutes(db *mongo.Database) *mux.Router {
	router := mux.NewRouter()

	controllers.SetCustomerCollection(db)

	// Rutas principales
	router.HandleFunc("/customers/{id}", controllers.DeleteCustomer).Methods("DELETE")

	router.HandleFunc("/sync-update", controllers.SyncUpdateCustomer).Methods("POST") // 🔥 Usa `SyncUpdateCustomer`
	router.HandleFunc("/sync-create", controllers.SyncCreateCustomer).Methods("POST")
	router.HandleFunc("/sync-delete/{id}", controllers.SyncDeleteCustomer).Methods("DELETE")

	return router
}

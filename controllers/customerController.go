package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"ProjectAWSStore-DeleteCustomer/config"
	"ProjectAWSStore-DeleteCustomer/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// 📌 Variable global para la colección de clientes
var customerCollection *mongo.Collection

// 📌 **Inicializar la colección en DeleteCustomer**
func SetCustomerCollection(db *mongo.Database) {
	customerCollection = db.Collection("customers")
}

// 📌 **Eliminar un cliente en DeleteCustomer**
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// 📌 Convertir `ID` de string a `primitive.ObjectID`
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "❌ ID inválido", http.StatusBadRequest)
		return
	}

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "❌ Database not initialized", http.StatusInternalServerError)
		return
	}

	// 🔍 Buscar si el cliente existe antes de eliminarlo
	var existingCustomer models.Customer
	err = customerCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existingCustomer)
	if err != nil {
		http.Error(w, "⚠️ Cliente no encontrado en DeleteCustomer", http.StatusNotFound)
		return
	}

	// 🗑️ Eliminar el cliente
	_, err = customerCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "❌ Error al eliminar cliente", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente eliminado correctamente:", existingCustomer.Email)

	// 🔄 **Sincronizar con ReadCustomer y CreateCustomer**
	go syncDeleteWithMicroservices(id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cliente eliminado correctamente"})
}

// 📌 **Notificar la eliminación a los otros microservicios**
func syncDeleteWithMicroservices(customerID string) {
	services := []string{
		os.Getenv("READ_CUSTOMER_SERVICE") + "/sync-delete",
		os.Getenv("CREATE_CUSTOMER_SERVICE") + "/sync-delete",
		os.Getenv("UPDATE_CUSTOMER_SERVICE") + "/sync-delete",
	}

	for _, service := range services {
		if service == "" {
			fmt.Println("⚠️ Servicio no definido, omitiendo...")
			continue
		}

		url := service + "/" + customerID
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			fmt.Println("❌ Error creando solicitud HTTP:", err)
			continue
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("❌ Error enviando solicitud a", url, ":", err)
			continue
		}

		fmt.Println("✅ Sincronización de eliminación exitosa con:", url, " Status:", resp.Status)
		resp.Body.Close()
	}
}

// 📌 **Recibir solicitudes de actualización desde otros microservicios**
// 📌 **Recibir solicitudes de actualización desde otros microservicios**
func SyncUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var updatedCustomer models.Customer

	// 📌 Decodificar JSON recibido
	err := json.NewDecoder(r.Body).Decode(&updatedCustomer)
	if err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de sincronización de actualización para:", updatedCustomer.Email)

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "❌ Database not initialized", http.StatusInternalServerError)
		return
	}

	// 📌 **Usar directamente `updatedCustomer.ID` como filtro** (sin `ObjectIDFromHex`)
	filter := bson.M{"_id": updatedCustomer.ID}
	update := bson.M{"$set": updatedCustomer}

	result, err := customerCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "❌ Error al sincronizar actualización en DeleteCustomer", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		fmt.Println("⚠️ Cliente no encontrado en la base de datos durante sincronización de actualización.")
		http.Error(w, "⚠️ Cliente no encontrado en la base de datos durante sincronización de actualización.", http.StatusNotFound)
		return
	}

	fmt.Println("✅ Cliente sincronizado correctamente en DeleteCustomer:", updatedCustomer.Email)
	w.WriteHeader(http.StatusOK)
}

// 📌 **Recibir solicitudes de eliminación desde otros microservicios**
func SyncDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	var deletedData struct {
		ID string `json:"id"`
	}

	// 📌 Decodificar JSON
	err := json.NewDecoder(r.Body).Decode(&deletedData)
	if err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de sincronización de eliminación para:", deletedData.ID)

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "❌ Database not initialized", http.StatusInternalServerError)
		return
	}

	// 📌 Convertir `ID` a `primitive.ObjectID`
	objID, err := primitive.ObjectIDFromHex(deletedData.ID)
	if err != nil {
		fmt.Println("⚠️ ID inválido en sincronización de eliminación:", deletedData.ID)
		http.Error(w, "⚠️ ID inválido en sincronización", http.StatusBadRequest)
		return
	}

	// 📌 Intentar eliminar el cliente en la base de datos
	result, err := customerCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "❌ Error al eliminar cliente en sincronización", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		fmt.Println("⚠️ Cliente no encontrado en la base de datos durante sincronización de eliminación.")
		http.Error(w, "⚠️ Cliente no encontrado en la base de datos durante sincronización de eliminación.", http.StatusNotFound)
		return
	}

	fmt.Println("✅ Cliente eliminado correctamente en DeleteCustomer:", deletedData.ID)
	w.WriteHeader(http.StatusOK)
}

// 📌 **Sincronizar clientes desde `CreateCustomer`**
func SyncCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de sincronización desde CreateCustomer:", customer.Email)

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// ✅ Verificar si el cliente ya existe en UpdateCustomer
	var existingCustomer models.Customer
	err = customerCollection.FindOne(context.TODO(), bson.M{"email": customer.Email}).Decode(&existingCustomer)
	if err == nil {
		fmt.Println("⚠️ Cliente ya existe en UpdateCustomer:", customer.Email)
		w.WriteHeader(http.StatusOK)
		return
	}

	// ✅ Insertar nuevo cliente
	_, err = customerCollection.InsertOne(context.TODO(), customer)
	if err != nil {
		http.Error(w, "❌ Error al sincronizar cliente", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente sincronizado correctamente en UpdateCustomer:", customer.Email)
	w.WriteHeader(http.StatusCreated)
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Srujankm12/CustomerPoBackend/internal/models"
	"github.com/Srujankm12/CustomerPoBackend/pkg/utils"
	"github.com/Srujankm12/CustomerPoBackend/repository"
	"github.com/gorilla/mux"
)

type CustomerPoHandler struct {
	customerRepo *repository.CustomerPoRepository
}

func NewCustomerPoHandler(customerRepo models.CustomerPoInterface) *CustomerPoHandler {
	return &CustomerPoHandler{
		customerRepo: customerRepo.(*repository.CustomerPoRepository),
	}
}
func (c *CustomerPoHandler) FetchDropDown(w http.ResponseWriter, r *http.Request) {
	customerList, err := c.customerRepo.FetchDropDown()
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Internal server error"})
		return
	}
	if len(customerList) == 0 {
		log.Printf("No data found")
		w.WriteHeader(http.StatusNotFound)
		utils.Encode(w, map[string]string{"message": "No data found"})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, customerList)

}
func (c *CustomerPoHandler) SubmitFormCustomerPoData(w http.ResponseWriter, r *http.Request) {
	var data models.CustomerPo
	err := utils.Decode(r, &data)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid request body"})
		return
	}

	err = c.customerRepo.SubmitFormCustomerPoData(data)
	if err != nil {
		log.Printf("Failed to submit customer PO data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Failed to submit data"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.Encode(w, map[string]string{"message": "Data submitted successfully"})
}

func (c *CustomerPoHandler) FetchCustomerPoData(w http.ResponseWriter, r *http.Request) {
	customerPoList, err := c.customerRepo.FetchCustomerPoData(r)
	if err != nil {
		log.Printf("Failed to fetch customer PO data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Internal server error"})
		return
	}
	if len(customerPoList) == 0 {
		log.Printf("No data found")
		w.WriteHeader(http.StatusNotFound)
		utils.Encode(w, map[string]string{"message": "No data found"})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, customerPoList)
}

func (c *CustomerPoHandler) UpdateCustomerPoData(w http.ResponseWriter, r *http.Request) {
	var data models.CustomerPo
	err := utils.Decode(r, &data)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid request body"})
		return
	}
	err = c.customerRepo.UpdateCustomerPoData(data)
	if err != nil {
		log.Printf("Failed to update customer PO data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Failed to update data"})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": "Data updated successfully"})
}

func (c *CustomerPoHandler) DeleteCustomerPoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = c.customerRepo.DeleteCustomerPo(id)
	if err != nil {
		log.Printf("Error deleting record: %v", err)
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Record deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

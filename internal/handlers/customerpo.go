package handlers

import (
	"log"
	"net/http"

	"github.com/Srujankm12/CustomerPoBackend/internal/models"
	"github.com/Srujankm12/CustomerPoBackend/pkg/utils"
)

type CustomerPoHandler struct {
	customerRepo models.CustomerPoInterface
}

func NewCustomerPoController(customerRepo models.CustomerPoInterface) *CustomerPoHandler {
	return &CustomerPoHandler{
		customerRepo: customerRepo,
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

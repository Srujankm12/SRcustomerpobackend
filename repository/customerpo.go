package repository

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Srujankm12/CustomerPoBackend/internal/models"
	"github.com/Srujankm12/CustomerPoBackend/pkg/database"
)

type CustomerPoRepository struct {
	db *sql.DB
}

func NewCustomerPoRepository(db *sql.DB) *CustomerPoRepository {
	return &CustomerPoRepository{
		db: db,
	}
}
func (c *CustomerPoRepository) FetchDropDown() ([]models.CustomerPoDropDown, error) {
	query := database.NewQuery(c.db)
	res, err := query.FetchDropDown(100, 0)
	if err != nil {
		log.Printf("Database query failed: %v", err)
		return nil, err
	}
	if len(res) == 0 {
		log.Println("No data found in FetchDropDown query")
		return nil, sql.ErrNoRows
	}
	log.Println("Successfully fetched dropdown data")
	return res, nil
}

func (c *CustomerPoRepository) SubmitFormCustomerPoData(data models.CustomerPo) error {
	query := database.NewQuery(c.db)
	err := query.SubmitFormCustomerPoData(data)
	if err != nil {
		log.Printf("Failed to submit customer PO data: %v", err)
		return err
	}
	return nil
}

func (c *CustomerPoRepository) FetchCustomerPoData(r *http.Request) ([]models.CustomerPo, error) {
	query := database.NewQuery(c.db)
	res, err := query.FetchCustomerPoData()
	if err != nil {
		log.Printf("Failed to fetch customer PO data: %v", err)
		return nil, err
	}
	return res, nil
}

func (c *CustomerPoRepository) UpdateCustomerPoData(data models.CustomerPo) error {
	query := database.NewQuery(c.db)
	err := query.UpdateCustomerPoData(data)
	if err != nil {
		log.Printf("Failed to update customer PO data: %v", err)
		return err
	}
	return nil
}

func (c *CustomerPoRepository) DeleteCustomerPo(id int) error {
	tx, err := c.db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}

	// Ensure rollback happens if there's an error
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Delete record based on the given ID
	_, err = tx.Exec("DELETE FROM customerposubmitteddata WHERE id = $1", id)
	if err != nil {
		log.Printf("Failed to delete record with id %d: %v", id, err)
		return err
	}

	return nil
}

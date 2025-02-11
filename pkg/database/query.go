package database

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/Srujankm12/CustomerPoBackend/internal/models"
)

type Query struct {
	db   *sql.DB
	Time *time.Location
}

func NewQuery(db *sql.DB) *Query {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatalf("Failed to load time zone: %v", err)
	}

	return &Query{
		db:   db,
		Time: loc,
	}
}

func (q *Query) CreateTables() error {
	tx, err := q.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS customername (
			customer_name VARCHAR(255) PRIMARY KEY
		)`,
		`CREATE TABLE IF NOT EXISTS bsno (
			bs_number VARCHAR(50) PRIMARY KEY
		)`,
		`CREATE TABLE IF NOT EXISTS unit (	
			unit_name VARCHAR(100) PRIMARY KEY
		)`,
		`CREATE TABLE IF NOT EXISTS postatusdd (
			status VARCHAR(100) PRIMARY KEY
		)`,
		`CREATE TABLE IF NOT EXISTS concernsonorder (
			concern VARCHAR(100) PRIMARY KEY
		)`,
		`CREATE TABLE IF NOT EXISTS customerposubmitteddata (
			timestamp VARCHAR(50) NOT NULL,
			sra_engineer_name VARCHAR(255),
			supplier VARCHAR(255),
			customer_name VARCHAR(255),
			bs_number VARCHAR(50),
			customer_po_no VARCHAR(100),
			po_date VARCHAR(50),
			part_code VARCHAR(100),
			qty INT,
			unit VARCHAR(100),
			total_value INT,
			po_status_dd VARCHAR(100),
			concerns_on_order VARCHAR(255),
			billable_sch_value INT,
			deli_sch_as_per_customer_po VARCHAR(255),
			customer_clearence_for_billing VARCHAR(255),
			reserved_qty_from_stock INT,
			required_qty_to_order INT,
			pending_qty_against_po INT,
			material_due_qty INT,
			so_number VARCHAR(100),
			mei_po_no VARCHAR(100),
			po_status_f VARCHAR(100),
			pending_value_against_po INT,
			pending_order_value INT,
			reserved_qty_stock_value INT,
			month_of_delivery_scheduled VARCHAR(50),
			category VARCHAR(100),
			PRIMARY KEY (customer_po_no)
		)`,
	}

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			log.Printf("Failed to execute query: %s", query)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	log.Println("All tables created successfully.")
	return nil
}

func (q *Query) FetchDropDown(limit, offset int) ([]models.CustomerPoDropDown, error) {
	var customerList []models.CustomerPoDropDown

	// 🚀 FIXED: Removed `CROSS JOIN` and added LIMIT + OFFSET for pagination
	query := `SELECT c.customer_name, b.bs_number, u.unit_name, s.status, o.concern
			  FROM customername c
			  JOIN bsno b ON TRUE
			  JOIN unit u ON TRUE
			  JOIN postatusdd s ON TRUE
			  JOIN concernsonorder o ON TRUE
			  LIMIT $1 OFFSET $2;`

	rows, err := q.db.Query(query, limit, offset)
	if err != nil {
		log.Printf("Database query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var customer models.CustomerPoDropDown
		if err := rows.Scan(&customer.CustomerName, &customer.BSNO, &customer.Unit, &customer.POStatusDD, &customer.ConcernsOnOrder); err != nil {
			log.Printf("Row scan failed: %v", err)
			return nil, err
		}
		customerList = append(customerList, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Printf("Fetched %d dropdown records", len(customerList))
	return customerList, nil
}
func (q *Query) SubmitFormCustomerPoData(data models.CustomerPo) error {
	// Convert customer clearance to an integer
	customerClearance, err := strconv.Atoi(data.CustomerClearanceForBilling)
	if err != nil {
		log.Printf("Invalid customer_clearence_for_billing value: %v", err)
		return err
	}

	// Initialize computed values
	var billableSchValue, pendingValueAgainstPO, pendingOrderValue, reservedQtyStockValue float64
	var requiredQtyToOrder, pendingQtyAgainstPO, materialDueQty int
	var poStatus string

	if data.Quantity > 0 {
		unitPrice := float64(data.TotalValue) / float64(data.Quantity) // Calculate unit price safely
		billableSchValue = unitPrice * float64(customerClearance)
		requiredQtyToOrder = data.Quantity - customerClearance - data.ReservedQtyFromStock
		pendingQtyAgainstPO = data.Quantity - customerClearance
		materialDueQty = data.Quantity - customerClearance - data.ReservedQtyFromStock
		pendingValueAgainstPO = float64(data.TotalValue) - billableSchValue
		pendingOrderValue = unitPrice * float64(requiredQtyToOrder)
		reservedQtyStockValue = unitPrice * float64(data.ReservedQtyFromStock)

		if pendingQtyAgainstPO == 0 {
			poStatus = "Completed"
		} else {
			poStatus = "Pending"
		}
	} else {
		// Handle zero quantity case safely
		billableSchValue = 0
		requiredQtyToOrder = 0
		pendingQtyAgainstPO = 0
		materialDueQty = 0
		pendingValueAgainstPO = 0
		pendingOrderValue = 0
		reservedQtyStockValue = 0
		poStatus = "Pending"
	}

	// Insert into the database
	_, err = q.db.Exec(`
		INSERT INTO customerposubmitteddata (
			timestamp,
			sra_engineer_name,
			supplier,
			customer_name,
			bs_number,
			customer_po_no,
			po_date,
			part_code,
			qty,
			unit,
			total_value,
			po_status_dd,
			concerns_on_order,
			billable_sch_value,
			deli_sch_as_per_customer_po,
			customer_clearence_for_billing,
			reserved_qty_from_stock,
			required_qty_to_order,
			pending_qty_against_po,
			material_due_qty,
			so_number,
			mei_po_no,
			po_status_f,
			pending_value_against_po,
			pending_order_value,
			reserved_qty_stock_value,
			month_of_delivery_scheduled,
			category
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)`,
		data.Timestamp,
		data.SRAEngineerName,
		data.Supplier,
		data.CustomerName,
		data.BSNO,
		data.CustomerPoNo,
		data.PODate,
		data.PartCode,
		data.Quantity,
		data.Unit,
		data.TotalValue,
		poStatus, // Computed PO Status
		data.ConcernsOnOrder,
		billableSchValue, // Computed Billable Sch Value
		data.DeliSchAsPerCustomerPo,
		customerClearance, // Converted clearance
		data.ReservedQtyFromStock,
		requiredQtyToOrder,  // Computed Required Qty to Order
		pendingQtyAgainstPO, // Computed Pending Qty Against PO
		materialDueQty,      // Computed Material Due Qty
		data.SONumber,
		data.MEIPONO,
		data.POStatusF,
		pendingValueAgainstPO, // Computed Pending Value Against PO
		pendingOrderValue,     // Computed Pending Order Value
		reservedQtyStockValue, // Computed Reserved Qty Stock Value
		data.MonthOfDeliveryScheduled,
		data.Category,
	)

	if err != nil {
		log.Printf("Failed to insert data into customerposubmitteddata: %v", err)
		return err
	}

	log.Println("Customer PO data submitted successfully.")
	return nil
}

func (q *Query) FetchCustomerPoData() ([]models.CustomerPo, error) {
	var customerPoList []models.CustomerPo

	query := `
		SELECT timestamp, sra_engineer_name, supplier, customer_name, bs_number, 
		       customer_po_no, po_date, part_code, qty, unit, total_value, 
		       po_status_dd, concerns_on_order, billable_sch_value, 
		       deli_sch_as_per_customer_po, customer_clearence_for_billing, 
		       reserved_qty_from_stock, required_qty_to_order, pending_qty_against_po, 
		       material_due_qty, so_number, mei_po_no, po_status_f, 
		       pending_value_against_po, pending_order_value, reserved_qty_stock_value, 
		       month_of_delivery_scheduled, category
		FROM customerposubmitteddata;
	`

	rows, err := q.db.Query(query)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var customerPo models.CustomerPo
		err := rows.Scan(
			&customerPo.Timestamp,
			&customerPo.SRAEngineerName,
			&customerPo.Supplier,
			&customerPo.CustomerName,
			&customerPo.BSNO,
			&customerPo.CustomerPoNo,
			&customerPo.PODate,
			&customerPo.PartCode,
			&customerPo.Quantity,
			&customerPo.Unit,
			&customerPo.TotalValue,
			&customerPo.POStatusDD,
			&customerPo.ConcernsOnOrder,
			&customerPo.BillableSchValue,
			&customerPo.DeliSchAsPerCustomerPo,
			&customerPo.CustomerClearanceForBilling,
			&customerPo.ReservedQtyFromStock,
			&customerPo.RequiredQtyToOrder,
			&customerPo.PendingQtyAgainstPO,
			&customerPo.MaterialDueQty,
			&customerPo.SONumber,
			&customerPo.MEIPONO,
			&customerPo.POStatusF,
			&customerPo.PendingValueAgainstPO,
			&customerPo.PendingOrderValue,
			&customerPo.ReservedQtyStockValue,
			&customerPo.MonthOfDeliveryScheduled,
			&customerPo.Category,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		customerPoList = append(customerPoList, customerPo)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	if len(customerPoList) == 0 {
		log.Println("No records found in customerposubmitteddata")
		return nil, nil // No error, just an empty result
	}

	log.Printf("Fetched %d records from customerposubmitteddata", len(customerPoList))
	return customerPoList, nil
}

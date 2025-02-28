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
			customer_name VARCHAR(255) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS bsno (
			bs_number VARCHAR(50) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS unit (	
			unit_name VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS postatusdd (
			status VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS concernsonorder (
			concern VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS customerposubmitteddata (
			id SERIAL PRIMARY KEY,
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
			total_value FLOAT, 
			po_status_dd VARCHAR(100),
			concerns_on_order VARCHAR(255),
			billable_sch_value FLOAT, 
			deli_sch_as_per_customer_po VARCHAR(255),
			customer_clearence_for_billing INT, 
			reserved_qty_from_stock INT,
			required_qty_to_order INT,
			pending_qty_against_po INT,
			material_due_qty INT,
			so_number VARCHAR(100),
			mei_po_no VARCHAR(100),
			po_status_final VARCHAR(100),
			pending_value_against_po FLOAT, 
			pending_order_value FLOAT,
			reserved_qty_stock_value FLOAT,
			month_of_delivery_scheduled VARCHAR(50),
			category VARCHAR(100)
			)`,
		`CREATE TABLE IF NOT EXISTS bengineername(
			b_engineer_name VARCHAR(255) NOT NULL UNIQUE
			) `,
		`CREATE TABLE IF NOT EXISTS bsupplier (
				b_supplier_name VARCHAR(255) NOT NULL UNIQUE
			)`,
		`CREATE TABLE IF NOT EXISTS bcustomer(
			b_customer_name VARCHAR(255) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS bunit(
			b_unit VARCHAR(100) NOT NULL UNIQUE
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

func (q *Query) FetchDropDown() ([]models.CustomerPoDropDown, error) {
	var dropdowns []models.CustomerPoDropDown

	rows, err := q.db.Query(`
		SELECT c.customer_name, b.bs_number, u.unit_name, s.status, o.concern
		FROM customername c
		CROSS JOIN bsno b
		CROSS JOIN unit u
		CROSS JOIN postatusdd s
		CROSS JOIN concernsonorder o;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dropdown models.CustomerPoDropDown
		if err := rows.Scan(&dropdown.CustomerName, &dropdown.BSNO, &dropdown.Unit, &dropdown.POStatusDD, &dropdown.ConcernsOnOrder); err != nil {
			return nil, err
		}
		dropdowns = append(dropdowns, dropdown)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dropdowns, nil
}

func (q *Query) SubmitFormCustomerPoData(data models.CustomerPo) error {
	customerClearance, err := strconv.Atoi(data.CustomerClearanceForBilling)
	if err != nil {
		log.Printf("Invalid customer_clearence_for_billing value: %v", err)
		return err
	}

	var billableSchValue, pendingValueAgainstPO, pendingOrderValue, reservedQtyStockValue float64
	var requiredQtyToOrder, pendingQtyAgainstPO, materialDueQty int
	var POStatusF string

	if data.Quantity > 0 {
		unitPrice := float64(data.TotalValue) / float64(data.Quantity)
		billableSchValue = unitPrice * float64(customerClearance)
		requiredQtyToOrder = data.Quantity - customerClearance - data.ReservedQtyFromStock
		pendingQtyAgainstPO = data.Quantity - customerClearance
		materialDueQty = data.Quantity - customerClearance - data.ReservedQtyFromStock
		pendingValueAgainstPO = float64(data.TotalValue) - billableSchValue
		pendingOrderValue = unitPrice * float64(requiredQtyToOrder)
		reservedQtyStockValue = unitPrice * float64(data.ReservedQtyFromStock)

		if pendingQtyAgainstPO == 0 {
			POStatusF = "Completed"
		} else {
			POStatusF = "Pending"
		}
	} else {

		billableSchValue = 0
		requiredQtyToOrder = 0
		pendingQtyAgainstPO = 0
		materialDueQty = 0
		pendingValueAgainstPO = 0
		pendingOrderValue = 0
		reservedQtyStockValue = 0
		POStatusF = "Pending"
	}

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
			po_status_final,
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
		data.POStatusDD,
		data.ConcernsOnOrder,
		billableSchValue,
		data.DeliSchAsPerCustomerPo,
		customerClearance,
		data.ReservedQtyFromStock,
		requiredQtyToOrder,
		pendingQtyAgainstPO,
		materialDueQty,
		data.SONumber,
		data.MEIPONO,
		POStatusF,
		pendingValueAgainstPO,
		pendingOrderValue,
		reservedQtyStockValue,
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

func (q *Query) UpdateCustomerPoData(data models.CustomerPo) error {
	customerClearance, err := strconv.Atoi(data.CustomerClearanceForBilling)
	if err != nil {
		log.Printf("Invalid customer_clearence_for_billing value: %v", err)
		return err
	}

	var billableSchValue, pendingValueAgainstPO, pendingOrderValue, reservedQtyStockValue float64
	var requiredQtyToOrder, pendingQtyAgainstPO, materialDueQty int
	var POStatusF string

	if data.Quantity > 0 {
		unitPrice := float64(data.TotalValue) / float64(data.Quantity)
		billableSchValue = unitPrice * float64(customerClearance)
		requiredQtyToOrder = data.Quantity - customerClearance - data.ReservedQtyFromStock
		pendingQtyAgainstPO = data.Quantity - customerClearance
		materialDueQty = data.Quantity - customerClearance - data.ReservedQtyFromStock
		pendingValueAgainstPO = float64(data.TotalValue) - billableSchValue
		pendingOrderValue = unitPrice * float64(requiredQtyToOrder)
		reservedQtyStockValue = unitPrice * float64(data.ReservedQtyFromStock)

		if pendingQtyAgainstPO == 0 {
			POStatusF = "Completed"
		} else {
			POStatusF = "Pending"
		}
	} else {

		billableSchValue = 0
		requiredQtyToOrder = 0
		pendingQtyAgainstPO = 0
		materialDueQty = 0
		pendingValueAgainstPO = 0
		pendingOrderValue = 0
		reservedQtyStockValue = 0
		POStatusF = "Pending"
	}

	_, err = q.db.Exec(`
		UPDATE customerposubmitteddata SET
			timestamp = $1,
			sra_engineer_name = $2,
			supplier = $3,
			customer_name = $4,
			bs_number = $5,
			customer_po_no = $6,
			po_date = $7,
			part_code = $8,
			qty = $9,
			unit = $10,
			total_value = $11,
			po_status_dd = $12,
			concerns_on_order = $13,
			billable_sch_value = $14,
			deli_sch_as_per_customer_po = $15,
			customer_clearence_for_billing = $16,
			reserved_qty_from_stock = $17,
			required_qty_to_order = $18,
			pending_qty_against_po = $19,
			material_due_qty = $20,
			so_number = $21,
			mei_po_no = $22,
			po_status_final = $23,
			pending_value_against_po = $24,
			pending_order_value = $25,
			reserved_qty_stock_value = $26,
			month_of_delivery_scheduled = $27,
			category = $28
		WHERE id = $29`,
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
		data.POStatusDD,
		data.ConcernsOnOrder,
		billableSchValue,
		data.DeliSchAsPerCustomerPo,
		customerClearance,
		data.ReservedQtyFromStock,
		requiredQtyToOrder,
		pendingQtyAgainstPO,
		materialDueQty,
		data.SONumber,
		data.MEIPONO,
		POStatusF,
		pendingValueAgainstPO,
		pendingOrderValue,
		reservedQtyStockValue,
		data.MonthOfDeliveryScheduled,
		data.Category,
		data.ID,
	)

	if err != nil {
		log.Printf("Failed to update data for ID %d: %v", data.ID, err)
		return err
	}

	log.Printf("Customer PO data updated successfully for ID %d.", data.ID)
	return nil
}

func (q *Query) FetchCustomerPoData() ([]models.CustomerPo, error) {
	var customerPoList []models.CustomerPo

	query := `
		SELECT id, timestamp, sra_engineer_name, supplier, customer_name, bs_number, 
		       customer_po_no, po_date, part_code, qty, unit, total_value, 
		       po_status_dd, concerns_on_order, billable_sch_value, 
		       deli_sch_as_per_customer_po, customer_clearence_for_billing, 
		       reserved_qty_from_stock, required_qty_to_order, pending_qty_against_po, 
		       material_due_qty, so_number, mei_po_no, po_status_final, 
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
			&customerPo.ID,
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
		return nil, nil
	}

	log.Printf("Fetched %d records from customerposubmitteddata", len(customerPoList))
	return customerPoList, nil
}

func (q *Query) DeleteCustomerPo(id int) error {
	tx, err := q.db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.Exec("DELETE FROM customerposubmitteddata WHERE id = $1", id)
	if err != nil {
		log.Printf("Failed to delete record with id %d: %v", id, err)
		return err
	}

	return nil
}

//billing sheet

// func (q *Query) FetchBillingSheetDropDown() ([]models.BillingSheetDropDown, error) {
// 	var billingSheetDropDown []models.BillingSheetDropDown
// 	rows, err := q.db.Query(`SELECT e.b_engineer_name,b.b_supplier_name,a.b_customer_name,u.b_unit
// 	FROM b_engineer e
// }

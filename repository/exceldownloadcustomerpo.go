package repository

import (
	"database/sql"
	"fmt"

	"github.com/Srujankm12/CustomerPoBackend/internal/models"
	"github.com/xuri/excelize/v2"
)

type ExcelDownloadCPO struct {
	db *sql.DB
}

func NewExcelDownloadCPO(db *sql.DB) *ExcelDownloadCPO {
	return &ExcelDownloadCPO{db: db}
}

func (e *ExcelDownloadCPO) FetchExcelCPO() ([]models.ExcelDownloadCPO, error) {
	var data []models.ExcelDownloadCPO

	if e.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	rows, err := e.db.Query("SELECT * FROM customerposubmitteddata")
	if err != nil {
		fmt.Println("Database query error:", err)
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var excelData models.ExcelDownloadCPO
		if err := rows.Scan(&excelData.ID, &excelData.Timestamp, &excelData.SRAEngineerName, &excelData.Supplier, &excelData.CustomerName, &excelData.BSNO, &excelData.CustomerPoNo, &excelData.PODate, &excelData.PartCode, &excelData.Quantity, &excelData.Unit, &excelData.TotalValue, &excelData.POStatusDD, &excelData.ConcernsOnOrder, &excelData.BillableSchValue, &excelData.DeliSchAsPerCustomerPo, &excelData.CustomerClearanceForBilling, &excelData.ReservedQtyFromStock, &excelData.RequiredQtyToOrder, &excelData.PendingQtyAgainstPO, &excelData.MaterialDueQty, &excelData.SONumber, &excelData.MEIPONO, &excelData.POStatusF, &excelData.PendingValueAgainstPO, &excelData.PendingOrderValue, &excelData.ReservedQtyStockValue, &excelData.MonthOfDeliveryScheduled, &excelData.Category); err != nil {
			fmt.Println("Error scanning row:", err)
			fmt.Println("Error scanning row:", err)
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		data = append(data, excelData)
	}
	if len(data) == 0 {
		fmt.Println("no data found in customerposubmitteddata table")
		return []models.ExcelDownloadCPO{}, nil
	}
	fmt.Printf("Fetched %d records\n", len(data))
	return data, nil
}
func (e *ExcelDownloadCPO) CreateExcelDownloadCPO() (*excelize.File, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	data, err := e.FetchExcelCPO()
	if err != nil {
		return nil, err
	}
	sheetName := "CustomerPO"
	file.NewSheet(sheetName)

	headers := []string{
		"ID",
		"Timestamp",
		"SRA Engineer Name",
		"Supplier",
		"Customer Name",
		"BS NO",
		"Customer PO No",
		"PO Date",
		"Part Code",
		"Quantity",
		"Unit",
		"Total Value",
		"PO Status DD",
		"Concerns on Order",
		"Billable Scheduled Value",
		"Delivery Schedule as per Customer PO",
		"Customer Clearance for Billing",
		"Reserved Quantity from Stock",
		"Required Quantity to Order",
		"Pending Quantity Against PO",
		"Material Due Quantity",
		"SO Number",
		"MEI PO NO",
		"PO Status Final",
		"Pending Value Against PO",
		"Pending Order Value",
		"Reserved Quantity Stock Value",
		"Month of Delivery Scheduled",
		"Category",
	}
	for colIndex, header := range headers {
		cell, err := excelize.CoordinatesToCellName(colIndex+1, 1)
		if err != nil {
			return nil, fmt.Errorf("error getting cell name: %v", err)
		}
		file.SetCellValue(sheetName, cell, header)
	}
	for i, record := range data {
		row := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), record.ID)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), record.Timestamp)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), record.SRAEngineerName)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), record.Supplier)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), record.CustomerName)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), record.BSNO)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", row), record.CustomerPoNo)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", row), record.PODate)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", row), record.PartCode)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", row), record.Quantity)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", row), record.Unit)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", row), record.TotalValue)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", row), record.POStatusDD)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", row), record.ConcernsOnOrder)
		file.SetCellValue(sheetName, fmt.Sprintf("O%d", row), record.BillableSchValue)
		file.SetCellValue(sheetName, fmt.Sprintf("P%d", row), record.DeliSchAsPerCustomerPo)
		file.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), record.CustomerClearanceForBilling)
		file.SetCellValue(sheetName, fmt.Sprintf("R%d", row), record.ReservedQtyFromStock)
		file.SetCellValue(sheetName, fmt.Sprintf("S%d", row), record.RequiredQtyToOrder)
		file.SetCellValue(sheetName, fmt.Sprintf("T%d", row), record.PendingQtyAgainstPO)
		file.SetCellValue(sheetName, fmt.Sprintf("U%d", row), record.MaterialDueQty)
		file.SetCellValue(sheetName, fmt.Sprintf("V%d", row), record.SONumber)
		file.SetCellValue(sheetName, fmt.Sprintf("W%d", row), record.MEIPONO)
		file.SetCellValue(sheetName, fmt.Sprintf("X%d", row), record.POStatusF)
		file.SetCellValue(sheetName, fmt.Sprintf("Y%d", row), record.PendingValueAgainstPO)
		file.SetCellValue(sheetName, fmt.Sprintf("Z%d", row), record.PendingOrderValue)
		file.SetCellValue(sheetName, fmt.Sprintf("AA%d", row), record.ReservedQtyStockValue)
		file.SetCellValue(sheetName, fmt.Sprintf("AB%d", row), record.MonthOfDeliveryScheduled)
		file.SetCellValue(sheetName, fmt.Sprintf("AC%d", row), record.Category)

	}
	return file, nil
}

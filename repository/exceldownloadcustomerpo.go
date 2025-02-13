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
		cell, _ := excelize.CoordinatesToCellName(colIndex+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}
	file.DeleteSheet("Sheet1")

	for i, record := range data {
		rowNum := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), record.ID)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), record.SRAEngineerName)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), record.Supplier)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), record.CustomerName)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), record.BSNO)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), record.CustomerPoNo)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), record.PODate)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), record.PartCode)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), record.Quantity)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), record.Unit)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), record.TotalValue)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), record.POStatusDD)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), record.ConcernsOnOrder)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", rowNum), record.BillableSchValue)
		file.SetCellValue(sheetName, fmt.Sprintf("O%d", rowNum), record.DeliSchAsPerCustomerPo)
		file.SetCellValue(sheetName, fmt.Sprintf("P%d", rowNum), record.CustomerClearanceForBilling)
		file.SetCellValue(sheetName, fmt.Sprintf("Q%d", rowNum), record.ReservedQtyFromStock)
		file.SetCellValue(sheetName, fmt.Sprintf("R%d", rowNum), record.RequiredQtyToOrder)
		file.SetCellValue(sheetName, fmt.Sprintf("S%d", rowNum), record.PendingQtyAgainstPO)
		file.SetCellValue(sheetName, fmt.Sprintf("T%d", rowNum), record.MaterialDueQty)
		file.SetCellValue(sheetName, fmt.Sprintf("U%d", rowNum), record.SONumber)
		file.SetCellValue(sheetName, fmt.Sprintf("V%d", rowNum), record.MEIPONO)
		file.SetCellValue(sheetName, fmt.Sprintf("W%d", rowNum), record.POStatusF)
		file.SetCellValue(sheetName, fmt.Sprintf("X%d", rowNum), record.PendingValueAgainstPO)
		file.SetCellValue(sheetName, fmt.Sprintf("Y%d", rowNum), record.PendingOrderValue)
		file.SetCellValue(sheetName, fmt.Sprintf("Z%d", rowNum), record.ReservedQtyStockValue)
		file.SetCellValue(sheetName, fmt.Sprintf("AA%d", rowNum), record.MonthOfDeliveryScheduled)
		file.SetCellValue(sheetName, fmt.Sprintf("AB%d", rowNum), record.Category)

	}
	return file, nil
}

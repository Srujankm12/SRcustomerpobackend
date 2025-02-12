package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Srujankm12/CustomerPoBackend/internal/models"
	"github.com/xuri/excelize/v2"
)

type ExcelDownloadCustomerPoHandler struct {
	customerRepo models.CustomerPoInterface
}

func NewExcelDownloadCustomerPoHandler(customerRepo models.CustomerPoInterface) *ExcelDownloadCustomerPoHandler {
	return &ExcelDownloadCustomerPoHandler{
		customerRepo: customerRepo,
	}
}
func (edc *ExcelDownloadCustomerPoHandler) DownloadCustomerPo(w http.ResponseWriter, r *http.Request) {
	data, err := edc.customerRepo.FetchCustomerPoData(r)
	if err != nil {
		http.Error(w, "Failed to fetch customer PO data", http.StatusInternalServerError)
		return
	}
	file := excelize.NewFile()
	sheetName := "CustomerPO"
	file.NewSheet(sheetName)

	headers := []string{
		"ID", "Timestamp", "SRA Engineer Name", "Supplier", "Customer Name", "BS Number",
		"Customer PO No", "PO Date", "Part Code", "Quantity", "Unit", "Total Value",
		"PO Status DD", "Concerns on Order", "Billable Sch Value", "Deli Sch as per Customer PO",
		"Customer Clearence for Billing", "Reserved Qty from Stock", "Required Qty to Order",
		"Pending Qty against PO", "Material Due Qty", "SO Number", "MEI PO No", "PO Status F",
		"Pending Value against PO", "Pending Order Value", "Reserved Qty Stock Value",
		"Month of Delivery Scheduled", "Category",
	}
	for colIndex, record := range headers {
		cell := fmt.Sprintf("%s%d", string(65+colIndex))
		file.SetCellValue(sheetName, cell, record)
	}

	for i, record := range data {
		rowNum := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), record.ID)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), record.Timestamp)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), record.SRAEngineerName)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), record.Supplier)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), record.CustomerName)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), record.BSNO)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), record.CustomerPoNo)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), record.PODate)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", rowNum), record.PartCode)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", rowNum), record.Quantity)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", rowNum), record.Unit)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", rowNum), record.TotalValue)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", rowNum), record.POStatusDD)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", rowNum), record.ConcernsOnOrder)
		file.SetCellValue(sheetName, fmt.Sprintf("O%d", rowNum), record.BillableSchValue)
		file.SetCellValue(sheetName, fmt.Sprintf("P%d", rowNum), record.DeliSchAsPerCustomerPo)
		file.SetCellValue(sheetName, fmt.Sprintf("Q%d", rowNum), record.CustomerClearanceForBilling)
		file.SetCellValue(sheetName, fmt.Sprintf("R%d", rowNum), record.ReservedQtyFromStock)
		file.SetCellValue(sheetName, fmt.Sprintf("S%d", rowNum), record.RequiredQtyToOrder)
		file.SetCellValue(sheetName, fmt.Sprintf("T%d", rowNum), record.PendingQtyAgainstPO)
		file.SetCellValue(sheetName, fmt.Sprintf("U%d", rowNum), record.MaterialDueQty)
		file.SetCellValue(sheetName, fmt.Sprintf("V%d", rowNum), record.SONumber)
		file.SetCellValue(sheetName, fmt.Sprintf("W%d", rowNum), record.MEIPONO)
		file.SetCellValue(sheetName, fmt.Sprintf("X%d", rowNum), record.POStatusF)
		file.SetCellValue(sheetName, fmt.Sprintf("Y%d", rowNum), record.PendingValueAgainstPO)
		file.SetCellValue(sheetName, fmt.Sprintf("Z%d", rowNum), record.PendingOrderValue)
		file.SetCellValue(sheetName, fmt.Sprintf("AA%d", rowNum), record.ReservedQtyStockValue)
		file.SetCellValue(sheetName, fmt.Sprintf("AB%d", rowNum), record.MonthOfDeliveryScheduled)
		file.SetCellValue(sheetName, fmt.Sprintf("AC%d", rowNum), record.Category)

	}
	tempDir := "/tmp"
	if os.Getenv("OS") == "Windows_NT" {
		tempDir = os.Getenv("TEMP")
	}
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create temporary directory", http.StatusInternalServerError)
		return
	}
	filepath := fmt.Sprintf("%s/customerpodata.xlsx", tempDir)
	if err := file.SaveAs(filepath); err != nil {
		http.Error(w, "Failed to save Excel file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=customerpodata.xlsx")
	http.ServeFile(w, r, filepath)

}

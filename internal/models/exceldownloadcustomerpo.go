package models

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type ExcelDownloadCPO struct {
	ID                          int     `json:"id"`
	Timestamp                   string  `json:"timestamp"`
	SRAEngineerName             string  `json:"sra_engineer_name"`
	Supplier                    string  `json:"supplier"`
	CustomerName                string  `json:"customer_name"`
	BSNO                        string  `json:"bs_no"`
	CustomerPoNo                string  `json:"customer_po_no"`
	PODate                      string  `json:"po_date"`
	PartCode                    string  `json:"part_code"`
	Quantity                    int     `json:"quantity"`
	Unit                        string  `json:"unit"`
	TotalValue                  float64 `json:"total_value"`
	POStatusDD                  string  `json:"po_status_dd"`
	ConcernsOnOrder             string  `json:"concerns_on_order"`
	BillableSchValue            float64 `json:"billable_scheduled_value"`
	DeliSchAsPerCustomerPo      string  `json:"delivery_schedule_as_per_customer_po"`
	CustomerClearanceForBilling string  `json:"customer_clearance_for_billing"`
	ReservedQtyFromStock        int     `json:"reserved_quantity_from_stock"`
	RequiredQtyToOrder          int     `json:"required_quantity_to_order"`
	PendingQtyAgainstPO         int     `json:"pending_quantity_against_po"`
	MaterialDueQty              int     `json:"material_due_quantity"`
	SONumber                    string  `json:"so_number"`
	MEIPONO                     string  `json:"mei_po_no"`
	POStatusF                   string  `json:"po_status_final"`
	PendingValueAgainstPO       float64 `json:"pending_value_against_po"`
	PendingOrderValue           float64 `json:"pending_order_value"`
	ReservedQtyStockValue       float64 `json:"reserved_quantity_stock_value"`
	MonthOfDeliveryScheduled    string  `json:"month_of_delivery_scheduled"`
	Category                    string  `json:"category"`
}

type ExcelDownloadCPOInterface interface {
	CreateExcelDownloadCPO(*io.ReadCloser) (*excelize.File, error)
}

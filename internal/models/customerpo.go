package models

import "net/http"

type CustomerPo struct {
	Timestamp                    string `json:"timestamp"`
	SRAEngineerName              string `json:"sra_engineer_name"`
	Supplier                     string `json:"supplier"`
	CustomerName                 string `json:"customer_name"`
	BSNO                         string `json:"bsno"`
	CustomerPoNo                 string `json:"customer_po_no"`
	PODate                       string `json:"po_date"`
	PartCode                     string `json:"part_code"`
	Quantity                     int    `json:"qty"`
	Unit                         string `json:"unit"`
	TotalValue                   int    `json:"total_value"`
	POStatusDD                   string `json:"po_status_dd"`
	ConcernsOnOrder              string `json:"concerns_on_order"`
	BillableSchValue             int    `json:"billable_sch_value"`
	DeliSchAsPerCustomerPo       string `json:"deli_sch_as_per_customer_po"`
	CustomerClearenceForBillingg string `json:"customer_clearence_for_billing"`
	ReservedQtyFromStock         int    `json:"reserved_qty_from_stock"`
	RequiredQtyToOrder           int    `json:"required_qty_to_order"`
	PendingQtyAgainstPO          int    `json:"pending_qty_against_po"`
	MaterialDueQty               int    `json:"material_due_qty"`
	SONumber                     string `json:"so_number"`
	MEIPONO                      string `json:"mei_po_no"`
	POStatusF                    string `json:"po_status_f"`
	PendingValueAgainstPO        int    `json:"pending_value_against_po"`
	PendingOrderValue            int    `json:"pending_order_value"`
	ReservedQtyStockValue        int    `json:"reserved_qty_stock_value"`
	MonthOfDeliveryScheduled     string `json:"month_of_delivery_scheduled"`
	Category                     string `json:"category"`
}
type CustomerPoDropDown struct {
	CustomerName    string `json:"customer_name"`
	BSNO            string `json:"bsno"`
	Unit            string `json:"unit"`
	POStatusDD      string `json:"po_status_dd"`
	ConcernsOnOrder string `json:"concerns_on_order"`
}

type CustomerPoInterface interface {
	FetchDropDown() ([]CustomerPoDropDown, error)
	SubmitFormCustomerPoData(customerPo CustomerPo) error
	FetchCustomerPoData(r *http.Request) ([]CustomerPo, error)
}

package models

type BillingSheet struct {
	ID                int    `json:"id"`
	Timestamp         string `json:"timestamp"`
	EngineerName      string `json:"b_engineer_name"`
	Supplier          string `json:"b_supplier"`
	BillNo            string `json:"bill_no"`
	Date              string `json:"date"`
	Customer          string `json:"b_customer"`
	CustomerPoNo      string `json:"customer_po_no"`
	PODate            string `json:"po_date"`
	ItemDescription   string `json:"item_description"`
	BilledQuantity    int    `json:"billed_quantity"`
	Unit              string `json:"b_unit"`
	NetValue          int    `json:"net_value"`
	CGST              int    `json:"cgst"`
	IGST              int    `json:"igst"`
	TotalTax          int    `json:"total_tax"`
	GrossValue        int    `json:"gross_value"`
	DispatchedThrough string `json:"dispatched_through"`
}

type BillingSheetDropDown struct {
	EngineerName string `json:"b_engineer_name"`
	Supplier     string `json:"b_supplier"`
	Customer     string `json:"b_customer"`
	Unit         string `json:"b_unit"`
}
type BillingSheetInterface interface {
	FetchDropDown() ([]BillingSheetDropDown, error)
}

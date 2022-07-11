package models

import "time"

// table schema for payement
type Payment struct {
	Payment_ID         uint      `json:"payment_id"`
	Digital_Payment    bool      `json:"digital_payment"`
	COD_Payement       bool      `json:"cod_payment"`
	Payment_Status     bool      `json:"payment_status"`
	Payment_Made_At    time.Time `json:"payment_made_at"`
	Payment_Updated_At time.Time `json:"payment_updated_at"`
}

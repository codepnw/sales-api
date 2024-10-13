package entities

import "time"

func TableCashier() string {
	return "cashiers"
}

type Cashier struct {
	Id        uint      `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name      string    `json:"name"`
	Passcode  string    `json:"passcode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CashierResponse struct {
	Id   uint   `json:"cashierId"`
	Name string `json:"name"`
}

type CashierListResponse struct {
	Cashiers []CashierResponse `json:"cashiers"`
	Meta     *QueryMeta        `json:"meta"`
}

package models

type Wallet struct {
	UUID           string `json:"uuid"`
	Amount         int    `json:"amount"`
	OpperationType string `json:"opperationType"`
}

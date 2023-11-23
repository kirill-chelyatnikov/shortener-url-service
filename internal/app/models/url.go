package models

type Link struct {
	ID            string
	BaseURL       string
	CorrelationID string
	Hash          string
	IsDeleted     bool
}

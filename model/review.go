package model

import "time"

type Review struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"service_id"`
	HouseholderID string    `json:"householder_id"`
	Rating        float64   `json:"rating"`
	Comments      string    `json:"comments"`
	ReviewDate    time.Time `json:"review_date"`
}

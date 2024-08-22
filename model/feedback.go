package model

type Feedback struct {
	ID            string  `json:"id"`
	HouseholderID string  `json:"householder_id"`
	ServiceID     string  `json:"service_id"`
	Comments      string  `json:"comments"`
	Rating        float64 `json:"rating"`
}

func NewFeedback(id, householderID, serviceID, comments string, rating float64) *Feedback {
	return &Feedback{
		ID:            id,
		HouseholderID: householderID,
		ServiceID:     serviceID,
		Comments:      comments,
		Rating:        rating,
	}
}

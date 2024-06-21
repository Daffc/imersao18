package service

type ReservationRequest struct {
	EventId    string   `json:"event_id"`
	Spots      []string `json:"spots"`
	TicketType string   `json:"ticket_type"`
	CardHash   string   `json:"card_hash"`
	Email      string   `json:"email"`
}

type ReservationResponse struct {
	Id         string `json:"id"`
	Email      string `json:"email"`
	Spot       string `json:"spot"`
	TicketType string `json:"ticket_kind"`
	Status     string `json:"status"`
	EventId    string `json:"event_id"`
}

type Partner interface {
	MakeReservation(req *ReservationRequest) ([]ReservationResponse, error)
}

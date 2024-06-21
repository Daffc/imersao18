package usecase

type EventDTO struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Location     string  `json:"location"`
	Organization string  `json:"organization"`
	Rating       string  `json:"rating"`
	Date         string  `json:"date"`
	ImageURL     string  `json:"image_url"`
	Capacity     int     `json:"capacity"`
	Price        float64 `json:"price"`
	PartnerId    int     `json:"partner_id"`
}

type SpotDTO struct {
	Id       string `json:"id"`
	EventId  string `json:"event_id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	TicketId string `json:"ticket_id"`
}

type TicketDTO struct {
	Id         string  `json:"id"`
	SpotId     string  `json:"spot_id"`
	TicketType string  `json:"ticket_type"`
	Price      float64 `json:"price"`
}

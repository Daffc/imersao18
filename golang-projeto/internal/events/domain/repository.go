package domain

type EventRepository interface {
	ListEvents() ([]Event, error)
	FindEventById(eventId string) (*Event, error)
	FindSpotsByEventId(eventId string) ([]*Spot, error)
	FindSpotByName(eventId, spotName string) (*Spot, error)
	// CreateEvent(event *Event) error
	CreateSpot(spot *Spot) error
	CreateTicket(ticket *Ticket) error
	ReserveSpot(spotId, ticketId string) error
}

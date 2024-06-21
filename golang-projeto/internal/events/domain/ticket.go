package domain

import (
	"errors"

	"github.com/google/uuid"
)

type TicketType string

const (
	TicketTypeHalf TicketType = "half"
	TicketTypeFull TicketType = "full"
)

type Ticket struct {
	Id         string
	EventId    string
	Spot       *Spot
	TicketType TicketType
	Price      float64
}

var (
	ErrTicketPriceLessThanZero = errors.New("ticker price must be greater than zero")
	ErrInvalidTicketType       = errors.New("invalid ticket type")
)

func IsValidTicketType(ticketType TicketType) bool {
	return ticketType == TicketTypeFull || ticketType == TicketTypeHalf
}

func (t *Ticket) CalculatePrice() {
	if t.TicketType == TicketTypeHalf {
		t.Price /= 2
	}
}

func (t *Ticket) Validate() error {
	if t.Price <= 0 {
		return ErrTicketPriceLessThanZero
	}

	return nil
}

func NewTicket(event *Event, spot *Spot, ticketType TicketType) (*Ticket, error) {
	if !IsValidTicketType(ticketType) {
		return nil, ErrInvalidTicketType
	}

	ticket := &Ticket{
		Id:         uuid.New().String(),
		EventId:    event.Id,
		Spot:       spot,
		TicketType: ticketType,
		Price:      event.Price,
	}

	ticket.CalculatePrice()
	if err := ticket.Validate(); err != nil {
		return nil, err
	}
	return ticket, nil
}

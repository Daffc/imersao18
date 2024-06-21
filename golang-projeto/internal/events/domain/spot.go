package domain

import (
	"errors"

	"github.com/google/uuid"
)

type SpotStatus string

const (
	SpotStatusAvailable SpotStatus = "available"
	SpotStatusSold      SpotStatus = "sold"
)

type Spot struct {
	Id       string
	EventId  string
	Name     string
	Status   SpotStatus
	TicketId string
}

var (
	ErrSpotNameRequired              = errors.New("invalid spot name")
	ErrSpotNameLessThanTwo           = errors.New("spot name must be at least 2 characters long")
	ErrInvalidSpotNameFirstCharacter = errors.New("spot name must start with a letter")
	ErrInvalidSpotNameLastCharacter  = errors.New("spot name must end with a number")
	ErrInvalidSpotNumber             = errors.New("invalid spot number")
	ErrSpotNotFound                  = errors.New("spot not found")
	ErrSpotAlreadyReserved           = errors.New("spot already reserved")
)

func (s *Spot) Validate() error {
	if s.Name == "" {
		return ErrSpotNameRequired
	}
	if len(s.Name) < 2 {
		return ErrSpotNameLessThanTwo
	}
	if s.Name[0] < 'A' || s.Name[0] > 'Z' {
		return ErrInvalidSpotNameFirstCharacter
	}
	if s.Name[1] < '0' || s.Name[1] > '9' {
		return ErrInvalidSpotNameLastCharacter
	}
	return nil
}

func NewSpot(event *Event, name string) (*Spot, error) {
	spot := &Spot{
		Id:      uuid.New().String(),
		EventId: event.Id,
		Name:    name,
		Status:  SpotStatusAvailable,
	}

	if err := spot.Validate(); err != nil {
		return nil, err
	}
	return spot, nil
}

func (s *Spot) Reserve(ticketId string) error {
	if s.Status == SpotStatusSold {
		return ErrSpotAlreadyReserved
	}

	s.Status = SpotStatusSold
	s.TicketId = ticketId
	return nil
}

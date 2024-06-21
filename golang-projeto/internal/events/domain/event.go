package domain

import (
	"errors"
	"time"
)

type Rating string

const (
	RatingLivre Rating = "L"
	Rating10    Rating = "L10"
	Rating12    Rating = "L12"
	Rating14    Rating = "L14"
	Rating16    Rating = "L16"
	Rating18    Rating = "L18"
)

type Event struct {
	Id           string
	Name         string
	Location     string
	Organization string
	Rating       Rating
	Date         time.Time
	ImageURL     string
	Capacity     int
	Price        float64
	PartnerId    int
	Spots        []Spot
	Tickets      []Ticket
}

var (
	ErrEventNameRequired          = errors.New("invalid event name")
	ErrEventInvalidDate           = errors.New("event date must by in the future")
	ErrEventCapacityLessEqualZero = errors.New("event capacity must be greater than zero")
	ErrEventPriceEqualZero        = errors.New("event price must be greater or equal to zero")
	ErrEventNotFound              = errors.New("event not found")
)

func (e *Event) Validate() error {
	if e.Name == "" {
		return ErrEventNameRequired
	}

	if e.Date.Before(time.Now()) {
		return ErrEventInvalidDate
	}

	if e.Capacity <= 0 {
		return ErrEventCapacityLessEqualZero
	}

	if e.Price < 0 {
		return ErrEventCapacityLessEqualZero
	}

	return nil
}

func (e *Event) AddSpot(name string) (*Spot, error) {
	spot, err := NewSpot(e, name)

	if err != nil {
		return nil, err
	}

	e.Spots = append(e.Spots, *spot)
	return spot, nil
}

package domain

import (
	"errors"
	"fmt"
)

type SpotService struct{}

func NewSpotService() *SpotService {
	return &SpotService{}
}

var (
	ErrorSpotServiceInvalidQuantity = errors.New("quantity must be greater than zero")
)

func (s *SpotService) GenerateSports(event *Event, quantity int) error {
	if quantity <= 0 {
		return ErrorSpotServiceInvalidQuantity
	}

	for i := range quantity {
		spotName := fmt.Sprintf("%c%d", 'A'+i/10, i%10+1)
		spot, err := NewSpot(event, spotName)
		if err != nil {
			return err
		}
		event.Spots = append(event.Spots, *spot)
	}

	return nil
}

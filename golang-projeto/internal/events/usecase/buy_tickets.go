package usecase

import (
	"github.com/daffc/imersao18/golang/internal/events/domain"
	"github.com/daffc/imersao18/golang/internal/events/infra/service"
)

type BuyTicketsInputDTO struct {
	EventId    string   `json:"event_id"`
	Spots      []string `json:"spots"`
	TicketType string   `json:"ticket_type"`
	CardHash   string   `json:"card_hash"`
	Email      string   `json:"email"`
}

type BuyTicketsOutputDTO struct {
	Tickets []TicketDTO `json:"tickets"`
}

type BuyTicketsUseCase struct {
	repo           domain.EventRepository
	partnerFactory service.PartnerFactory
}

func NewBuyTicketsUseCase(repo domain.EventRepository, partnerFactory service.PartnerFactory) *BuyTicketsUseCase {
	return &BuyTicketsUseCase{repo: repo, partnerFactory: partnerFactory}
}

func (uc *BuyTicketsUseCase) Execute(input BuyTicketsInputDTO) (*BuyTicketsOutputDTO, error) {

	event, err := uc.repo.FindEventById(input.EventId)
	if err != nil {
		return nil, err
	}

	req := &service.ReservationRequest{
		EventId:    input.EventId,
		Spots:      input.Spots,
		TicketType: input.TicketType,
		CardHash:   input.CardHash,
		Email:      input.Email,
	}

	partnerSerice, err := uc.partnerFactory.CreatePartner(event.PartnerId)
	if err != nil {
		return nil, err
	}

	reservationResponse, err := partnerSerice.MakeReservation(req)
	if err != nil {
		return nil, err
	}

	// Id
	// Email
	// Spot
	// TicketType
	// Status
	// EventId

	tickets := make([]domain.Ticket, len(reservationResponse))
	for i, reservation := range reservationResponse {
		// Recovering related spot
		spot, err := uc.repo.FindSpotByName(reservation.EventId, reservation.Spot)
		if err != nil {
			return nil, err
		}

		// Generating a new ticket
		ticket, err := domain.NewTicket(event, spot, domain.TicketType(input.TicketType))
		if err != nil {
			return nil, err
		}

		// Creating ticket (database)
		err = uc.repo.CreateTicket(ticket)
		if err != nil {
			return nil, err
		}

		// Reserving spot
		err = spot.Reserve(spot.TicketId)
		if err != nil {
			return nil, err
		}

		err = uc.repo.ReserveSpot(spot.Id, ticket.Id)
		if err != nil {
			return nil, err
		}

		tickets[i] = domain.Ticket{
			Id:         ticket.Id,
			Spot:       ticket.Spot,
			TicketType: ticket.TicketType,
			Price:      ticket.Price,
		}
	}

	ticketsDTO := make([]TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketsDTO[i] = TicketDTO{
			Id:         ticket.Id,
			SpotId:     ticket.Spot.Id,
			TicketType: string(ticket.TicketType),
			Price:      ticket.Price,
		}
	}

	return &BuyTicketsOutputDTO{Tickets: ticketsDTO}, nil

}

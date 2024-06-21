package usecase

import "github.com/daffc/imersao18/golang/internal/events/domain"

type ListSpotsInputDTO struct {
	EventId string
}

type ListSpotsOutputDTO struct {
	Event EventDTO  `json:"event"`
	Spots []SpotDTO `json:"spots"`
}

type ListSpotsUseCase struct {
	repo domain.EventRepository
}

func NewListSpotsUseCase(repo domain.EventRepository) *ListSpotsUseCase {
	return &ListSpotsUseCase{repo: repo}
}

func (uc *ListSpotsUseCase) Execute(input ListSpotsInputDTO) (*ListSpotsOutputDTO, error) {

	// Buscando dados em db.
	event, err := uc.repo.FindEventById(input.EventId)
	if err != nil {
		return nil, err
	}

	spots, err := uc.repo.FindSpotsByEventId(input.EventId)
	if err != nil {
		return nil, err
	}

	// Ajustando dados a DTO para serem entregues a cliente.
	eventDTO := EventDTO{
		Id:           event.Id,
		Name:         event.Name,
		Location:     event.Location,
		Organization: event.Organization,
		Rating:       string(event.Rating),
		Date:         event.Date.Format("0000-00-00 00:00:00"),
		ImageURL:     event.ImageURL,
		Capacity:     event.Capacity,
		Price:        event.Price,
		PartnerId:    event.PartnerId,
	}

	spotsDTO := make([]SpotDTO, len(spots))
	for i, spot := range spots {
		spotsDTO[i] = SpotDTO{
			Id:       spot.Id,
			EventId:  spot.EventId,
			Name:     spot.Name,
			Status:   string(spot.Status),
			TicketId: spot.TicketId,
		}
	}

	return &ListSpotsOutputDTO{Event: eventDTO, Spots: spotsDTO}, nil
}

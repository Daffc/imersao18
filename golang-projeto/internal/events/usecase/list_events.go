package usecase

import "github.com/daffc/imersao18/golang/internal/events/domain"

type ListEventsOutputDTO struct {
	Events []EventDTO
}

type ListEventsUseCase struct {
	repo domain.EventRepository
}

func NewListEvenetsUseCase(repo domain.EventRepository) *ListEventsUseCase {
	return &ListEventsUseCase{repo: repo}
}

func (uc *ListEventsUseCase) Execute() (*ListEventsOutputDTO, error) {

	// Buscando dados em db.
	events, err := uc.repo.ListEvents()
	if err != nil {
		return nil, err
	}

	eventsDTO := make([]EventDTO, len(events))

	// Ajustando dados a DTO para serem entregues a cliente.
	for i, event := range events {
		eventsDTO[i] = EventDTO{
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
	}

	return &ListEventsOutputDTO{Events: eventsDTO}, nil
}

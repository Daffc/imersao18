package usecase

import "github.com/daffc/imersao18/golang/internal/events/domain"

type GetEventInputDTO struct {
	Id string
}

type GetEventOutputDTO struct {
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

type GetEventsUseCase struct {
	repo domain.EventRepository
}

func NewGetEventUseCase(repo domain.EventRepository) *GetEventsUseCase {
	return &GetEventsUseCase{repo: repo}
}

func (uc *GetEventsUseCase) Execute(input GetEventInputDTO) (*GetEventOutputDTO, error) {

	// Buscando dados em db.
	event, err := uc.repo.FindEventById(input.Id)
	if err != nil {
		return nil, err
	}

	// Ajustando dados a DTO para serem entregues a cliente.
	eventDTO := GetEventOutputDTO{
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

	return &eventDTO, nil
}

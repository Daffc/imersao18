package http

import (
	"encoding/json"
	"net/http"

	"github.com/daffc/imersao18/golang/internal/events/usecase"
)

type EventsHandler struct {
	listEventsUseCase *usecase.ListEventsUseCase
	getEventsUseCase  *usecase.GetEventsUseCase
	listSpotsUseCase  *usecase.ListSpotsUseCase
	buyTicketsUseCase *usecase.BuyTicketsUseCase
}

func NewEventHandler(
	listEventsUseCase *usecase.ListEventsUseCase,
	getEventsUseCase *usecase.GetEventsUseCase,
	listSpotsUseCase *usecase.ListSpotsUseCase,
	buyTicketsUseCase *usecase.BuyTicketsUseCase,
) *EventsHandler {
	return &EventsHandler{
		listEventsUseCase: listEventsUseCase,
		getEventsUseCase:  getEventsUseCase,
		listSpotsUseCase:  listSpotsUseCase,
		buyTicketsUseCase: buyTicketsUseCase,
	}
}

func (h *EventsHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	output, err := h.listEventsUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *EventsHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventId := r.PathValue("eventId")
	input := usecase.GetEventInputDTO{Id: eventId}
	output, err := h.getEventsUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *EventsHandler) ListSpots(w http.ResponseWriter, r *http.Request) {
	eventId := r.PathValue("eventId")
	input := usecase.ListSpotsInputDTO{EventId: eventId}
	output, err := h.listSpotsUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (h *EventsHandler) BuyTickets(w http.ResponseWriter, r *http.Request) {
	var input usecase.BuyTicketsInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := h.buyTicketsUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

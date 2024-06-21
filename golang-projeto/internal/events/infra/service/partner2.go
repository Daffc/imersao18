package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Partner2 struct {
	BaseURL string
}

type Partner2ReservationRequest struct {
	Lugares      []string `json:"lugares"`
	TipoIngresso string   `json:"tipo_ingresso"`
	Email        string   `json:"email"`
	EventId      string   `json:"event_id"`
}

type Partner2ReservationResponse struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	Lugar        string `json:"lugar"`
	TipoIngresso string `json:"tipo_ingresso"`
	Status       string `json:"status"`
	EventId      string `json:"event_id"`
}

func (p *Partner2) MakeReservation(req *ReservationRequest) ([]ReservationResponse, error) {
	partnerRequest := Partner2ReservationRequest{
		Lugares:      req.Spots,
		TipoIngresso: req.TicketType,
		Email:        req.Email,
	}

	body, err := json.Marshal(partnerRequest)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/eventos/%s/reservar", p.BaseURL, req.EventId)
	httpRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	httpRespose, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpRespose.Body.Close()

	if httpRespose.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", httpRespose.StatusCode)
	}

	var partnerResponse []Partner2ReservationResponse
	if err := json.NewDecoder(httpRespose.Body).Decode(&partnerResponse); err != nil {
		return nil, err
	}

	responses := make([]ReservationResponse, len(partnerResponse))

	for i, r := range partnerResponse {
		responses[i] = ReservationResponse{
			Id:         r.Id,
			Email:      r.Email,
			TicketType: r.TipoIngresso,
			Spot:       r.Lugar,
			Status:     r.Status,
			EventId:    r.EventId,
		}
	}

	return responses, nil

}

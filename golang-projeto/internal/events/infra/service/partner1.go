package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Partner1 struct {
	BaseURL string
}

type Partner1ReservationRequest struct {
	Spots      []string `json:"spots"`
	TicketKind string   `json:"ticket_kind"`
	Email      string   `json:"email"`
	EventId    string   `json:"event_id"`
}

type Partner1ReservationResponse struct {
	Id         string `json:"id"`
	Email      string `json:"email"`
	Spot       string `json:"spot"`
	TicketKind string `json:"ticket_kind"`
	Status     string `json:"status"`
	EventId    string `json:"event_id"`
}

func (p *Partner1) MakeReservation(req *ReservationRequest) ([]ReservationResponse, error) {
	partnerRequest := Partner1ReservationRequest{
		Spots:      req.Spots,
		TicketKind: req.TicketType,
		Email:      req.Email,
	}

	body, err := json.Marshal(partnerRequest)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/events/%s/reserve", p.BaseURL, req.EventId)
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

	var partnerResponse []Partner1ReservationResponse
	if err := json.NewDecoder(httpRespose.Body).Decode(&partnerResponse); err != nil {
		return nil, err
	}

	responses := make([]ReservationResponse, len(partnerResponse))

	for i, r := range partnerResponse {
		responses[i] = ReservationResponse{
			Id:         r.Id,
			Email:      r.Email,
			TicketType: r.TicketKind,
			Spot:       r.Spot,
			Status:     r.Status,
			EventId:    r.EventId,
		}
	}

	return responses, nil

}

package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/daffc/imersao18/golang/internal/events/domain"
	_ "github.com/go-sql-driver/mysql"
)

type mysqlEventRepository struct {
	db *sql.DB
}

// NewMysqlEventRepository creates a new MySQL event repository.
func NewMysqlEventRepository(db *sql.DB) (domain.EventRepository, error) {
	return &mysqlEventRepository{db: db}, nil
}

// ListEvents returns all events with their associated spots and tickets.
func (r *mysqlEventRepository) ListEvents() ([]domain.Event, error) {
	query := `
		SELECT 
			e.id, e.name, e.location, e.organization, e.rating, e.date, e.image_url, e.capacity, e.price, e.partner_id,
			s.id, s.event_id, s.name, s.status, s.ticket_id,
			t.id, t.event_id, t.spot_id, t.ticket_type, t.price
		FROM events e
		LEFT JOIN spots s ON e.id = s.event_id
		LEFT JOIN tickets t ON s.id = t.spot_id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventMap := make(map[string]*domain.Event)
	spotMap := make(map[string]*domain.Spot)
	for rows.Next() {
		var eventId, eventName, eventLocation, eventOrganization, eventRating, eventImageURL, spotId, spotEventId, spotName, spotStatus, spotTicketId, ticketId, ticketEventId, ticketSpotId, ticketType sql.NullString
		var eventDate sql.NullString
		var eventCapacity int
		var eventPrice, ticketPrice sql.NullFloat64
		var partnerId sql.NullInt32

		err := rows.Scan(
			&eventId, &eventName, &eventLocation, &eventOrganization, &eventRating, &eventDate, &eventImageURL, &eventCapacity, &eventPrice, &partnerId,
			&spotId, &spotEventId, &spotName, &spotStatus, &spotTicketId,
			&ticketId, &ticketEventId, &ticketSpotId, &ticketType, &ticketPrice,
		)
		if err != nil {
			return nil, err
		}

		if !eventId.Valid || !eventName.Valid || !eventLocation.Valid || !eventOrganization.Valid || !eventRating.Valid || !eventDate.Valid || !eventImageURL.Valid || !eventPrice.Valid || !partnerId.Valid {
			continue
		}

		event, exists := eventMap[eventId.String]
		if !exists {
			eventDateParsed, err := time.Parse("2006-01-02 15:04:05", eventDate.String)
			if err != nil {
				return nil, err
			}
			event = &domain.Event{
				Id:           eventId.String,
				Name:         eventName.String,
				Location:     eventLocation.String,
				Organization: eventOrganization.String,
				Rating:       domain.Rating(eventRating.String),
				Date:         eventDateParsed,
				ImageURL:     eventImageURL.String,
				Capacity:     eventCapacity,
				Price:        eventPrice.Float64,
				PartnerId:    int(partnerId.Int32),
				Spots:        []domain.Spot{},
				Tickets:      []domain.Ticket{},
			}
			eventMap[eventId.String] = event
		}

		if spotId.Valid {
			spot, spotExists := spotMap[spotId.String]
			if !spotExists {
				spot = &domain.Spot{
					Id:       spotId.String,
					EventId:  spotEventId.String,
					Name:     spotName.String,
					Status:   domain.SpotStatus(spotStatus.String),
					TicketId: spotTicketId.String,
				}
				event.Spots = append(event.Spots, *spot)
				spotMap[spotId.String] = spot
			}

			if ticketId.Valid {
				ticket := domain.Ticket{
					Id:         ticketId.String,
					EventId:    ticketEventId.String,
					Spot:       spot,
					TicketType: domain.TicketType(ticketType.String),
					Price:      ticketPrice.Float64,
				}
				event.Tickets = append(event.Tickets, ticket)
			}
		}
	}

	events := make([]domain.Event, 0, len(eventMap))
	for _, event := range eventMap {
		events = append(events, *event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// FindEventById returns an event by its Id, including associated spots and tickets.
func (r *mysqlEventRepository) FindEventById(eventId string) (*domain.Event, error) {
	query := `
		SELECT 
			e.id, e.name, e.location, e.organization, e.rating, e.date, e.image_url, e.capacity, e.price, e.partner_id,
			s.id, s.event_id, s.name, s.status, s.ticket_id,
			t.id, t.event_id, t.spot_id, t.ticket_type, t.price
		FROM events e
		LEFT JOIN spots s ON e.id = s.event_id
		LEFT JOIN tickets t ON s.id = t.spot_id
		WHERE e.id = ?
	`
	rows, err := r.db.Query(query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var event *domain.Event
	for rows.Next() {
		var eventIdStr, eventName, eventLocation, eventOrganization, eventRating, eventImageURL, spotId, spotEventId, spotName, spotStatus, spotTicketId, ticketId, ticketEventId, ticketSpotId, ticketType sql.NullString
		var eventDate sql.NullString
		var eventCapacity int
		var eventPrice, ticketPrice sql.NullFloat64
		var partnerId sql.NullInt32

		err := rows.Scan(
			&eventIdStr, &eventName, &eventLocation, &eventOrganization, &eventRating, &eventDate, &eventImageURL, &eventCapacity, &eventPrice, &partnerId,
			&spotId, &spotEventId, &spotName, &spotStatus, &spotTicketId,
			&ticketId, &ticketEventId, &ticketSpotId, &ticketType, &ticketPrice,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrEventNotFound
			}
			return nil, err
		}

		if !eventIdStr.Valid || !eventName.Valid || !eventLocation.Valid || !eventOrganization.Valid || !eventRating.Valid || !eventDate.Valid || !eventImageURL.Valid || !eventPrice.Valid || !partnerId.Valid {
			continue
		}

		if event == nil {
			eventDateParsed, err := time.Parse("2006-01-02 15:04:05", eventDate.String)
			if err != nil {
				return nil, err
			}
			event = &domain.Event{
				Id:           eventIdStr.String,
				Name:         eventName.String,
				Location:     eventLocation.String,
				Organization: eventOrganization.String,
				Rating:       domain.Rating(eventRating.String),
				Date:         eventDateParsed,
				ImageURL:     eventImageURL.String,
				Capacity:     eventCapacity,
				Price:        eventPrice.Float64,
				PartnerId:    int(partnerId.Int32),
				Spots:        []domain.Spot{},
				Tickets:      []domain.Ticket{},
			}
		}

		if spotId.Valid {
			spot := domain.Spot{
				Id:       spotId.String,
				EventId:  spotEventId.String,
				Name:     spotName.String,
				Status:   domain.SpotStatus(spotStatus.String),
				TicketId: spotTicketId.String,
			}
			event.Spots = append(event.Spots, spot)

			if ticketId.Valid {
				ticket := domain.Ticket{
					Id:         ticketId.String,
					EventId:    ticketEventId.String,
					Spot:       &spot,
					TicketType: domain.TicketType(ticketType.String),
					Price:      ticketPrice.Float64,
				}
				event.Tickets = append(event.Tickets, ticket)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if event == nil {
		return nil, domain.ErrEventNotFound
	}

	return event, nil
}

// CreateEvent inserts a new event into the database.
func (r *mysqlEventRepository) CreateEvent(event *domain.Event) error {
	query := `
		INSERT INTO events (id, name, location, organization, rating, date, image_url, capacity, price, partner_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, event.Id, event.Name, event.Location, event.Organization, event.Rating, event.Date.Format("2006-01-02 15:04:05"), event.ImageURL, event.Capacity, event.Price, event.PartnerId)
	return err
}

// FindSpotById returns a spot by its Id, including the associated ticket (if any).
func (r *mysqlEventRepository) FindSpotById(spotId string) (*domain.Spot, error) {
	query := `
		SELECT
			s.id, s.event_id, s.name, s.status, s.ticket_id,
			t.id, t.event_id, t.spot_id, t.ticket_type, t.price
		FROM spots s
		LEFT JOIN tickets t ON s.id = t.spot_id
		WHERE s.id = ?
	`
	row := r.db.QueryRow(query, spotId)

	var spot domain.Spot
	var ticket domain.Ticket
	var ticketId, ticketEventId, ticketSpotId, ticketType sql.NullString
	var ticketPrice sql.NullFloat64

	err := row.Scan(
		&spot.Id, &spot.EventId, &spot.Name, &spot.Status, &spot.TicketId,
		&ticketId, &ticketEventId, &ticketSpotId, &ticketType, &ticketPrice,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSpotNotFound
		}
		return nil, err
	}

	if ticketId.Valid {
		ticket.Id = ticketId.String
		ticket.EventId = ticketEventId.String
		ticket.Spot = &spot
		ticket.TicketType = domain.TicketType(ticketType.String)
		ticket.Price = ticketPrice.Float64
		spot.TicketId = ticket.Id
	}

	return &spot, nil
}

// CreateSpot inserts a new spot into the database.
func (r *mysqlEventRepository) CreateSpot(spot *domain.Spot) error {
	query := `
		INSERT INTO spots (id, event_id, name, status, ticket_id)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, spot.Id, spot.EventId, spot.Name, spot.Status, spot.TicketId)
	return err
}

// CreateTicket inserts a new ticket into the database.
func (r *mysqlEventRepository) CreateTicket(ticket *domain.Ticket) error {
	query := `
		INSERT INTO tickets (id, event_id, spot_id, ticket_type, price)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, ticket.Id, ticket.EventId, ticket.Spot.Id, ticket.TicketType, ticket.Price)
	return err
}

// ReserveSpot updates a spot's status to reserved.
func (r *mysqlEventRepository) ReserveSpot(spotId, ticketId string) error {
	query := `
		UPDATE spots
		SET status = ?, ticket_id = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, domain.SpotStatusSold, ticketId, spotId)
	return err
}

// FindSpotsByEventId returns all spots for a given event Id.
func (r *mysqlEventRepository) FindSpotsByEventId(eventId string) ([]*domain.Spot, error) {
	query := `
		SELECT id, event_id, name, status, ticket_id
		FROM spots
		WHERE event_id = ?
	`
	rows, err := r.db.Query(query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spots []*domain.Spot
	for rows.Next() {
		var spot domain.Spot
		if err := rows.Scan(&spot.Id, &spot.EventId, &spot.Name, &spot.Status, &spot.TicketId); err != nil {
			return nil, err
		}
		spots = append(spots, &spot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return spots, nil
}

func (r *mysqlEventRepository) FindSpotByName(eventId, name string) (*domain.Spot, error) {
	query := `
		SELECT 
			s.id, s.event_id, s.name, s.status, s.ticket_id,
			t.id, t.event_id, t.spot_id, t.ticket_type, t.price
		FROM spots s
		LEFT JOIN tickets t ON s.id = t.spot_id
		WHERE s.event_id = ? AND s.name = ?
	`
	row := r.db.QueryRow(query, eventId, name)

	var spot domain.Spot
	var ticket domain.Ticket
	var ticketId, ticketEventId, ticketSpotId, ticketType sql.NullString
	var ticketPrice sql.NullFloat64

	err := row.Scan(
		&spot.Id, &spot.EventId, &spot.Name, &spot.Status, &spot.TicketId,
		&ticketId, &ticketEventId, &ticketSpotId, &ticketType, &ticketPrice,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSpotNotFound
		}
		return nil, err
	}

	if ticketId.Valid {
		ticket.Id = ticketId.String
		ticket.EventId = ticketEventId.String
		ticket.Spot = &spot
		ticket.TicketType = domain.TicketType(ticketType.String)
		ticket.Price = ticketPrice.Float64
		spot.TicketId = ticket.Id
	}

	return &spot, nil
}

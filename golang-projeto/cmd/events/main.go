package events

import (
	"database/sql"
	"net/http"

	"github.com/daffc/imersao18/golang/internal/events/infra/repository"
	"github.com/daffc/imersao18/golang/internal/events/infra/service"
	"github.com/daffc/imersao18/golang/internal/events/usecase"

	httpHandler "github.com/daffc/imersao18/golang/internal/events/infra/http"
)

func main() {
	// Conectando a banco de dados.
	db, err := sql.Open("mysql", "test_user:test_password@tcp(localhost:3306)/test_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventRepo, err := repository.NewMysqlEventRepository(db)
	if err != nil {
		panic(err)
	}

	// Definindo Partners
	partnerBaseURLs := map[int]string{
		1: "http://localjpst:9080/api1",
		2: "http://localjpst:9080/api2",
	}
	partnerFactory := service.NewPartnerfactory(partnerBaseURLs)

	// Definindo Rotas e HttpHandler
	listEventsUseCase := usecase.NewListEvenetsUseCase(eventRepo)
	getEventsUseCase := usecase.NewGetEventUseCase(eventRepo)
	listSpotsUseCase := usecase.NewListSpotsUseCase(eventRepo)
	buyTicketsUseCase := usecase.NewBuyTicketsUseCase(eventRepo, partnerFactory)

	eventsHandler := httpHandler.NewEventHandler(
		listEventsUseCase,
		getEventsUseCase,
		listSpotsUseCase,
		buyTicketsUseCase,
	)

	r := http.NewServeMux()
	r.HandleFunc("GET /events", eventsHandler.ListEvents)
	r.HandleFunc("GET /events/{eventId}", eventsHandler.GetEvent)
	r.HandleFunc("GET /events/{eventId}/spots", eventsHandler.ListSpots)
	r.HandleFunc("POST /checkout", eventsHandler.BuyTickets)

	http.ListenAndServe(":8080", r)
}

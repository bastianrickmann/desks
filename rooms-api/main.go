package main

import (
	"fmt"
	"github.com/ChristophBe/rooms/rooms-api/data"
	"github.com/ChristophBe/rooms/rooms-api/handlers"
	"github.com/ChristophBe/rooms/rooms-api/middlewares"
	"github.com/ChristophBe/rooms/rooms-api/util"
	"log"
	"net/http"
)
import "github.com/gorilla/mux"

func main() {

	urlPrefix := "/api/v1.0"
	router := mux.NewRouter()
	router.Path(urlPrefix + "/users/login").HandlerFunc(handlers.PostUsersLogin).Methods(http.MethodPost)
	router.Path(urlPrefix + "/users/me").Handler(withAuth(handlers.GetUsersMe)).Methods(http.MethodGet)
	router.Path(urlPrefix + "/rooms").Handler(withAuth(handlers.GetAllRooms)).Methods(http.MethodGet)

	router.Path(urlPrefix + "/users/{id}/bookings").Handler(withAuth(handlers.GetBookingsByUser)).Methods(http.MethodGet)
	router.Path(urlPrefix + "/bookings").Handler(withAuth(handlers.PostBooking)).Methods(http.MethodPost)

	if err := data.InitDatabase(); err != nil {
		log.Fatal("Failed initialize Database", err)
	}

	serverPort := util.GetIntEnvironmentVariable("SERVER_PORT", 8080)
	log.Printf("Starting Server and expose port %d\n", serverPort)

	httpRequestHandler := middlewares.CorsMiddleware(router)
	httpRequestHandler = middlewares.LoggingMiddleware(httpRequestHandler)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), httpRequestHandler); err != nil {
		log.Fatal("Failed to start Api", err)
	}
}

func withAuth(handlerFunc http.HandlerFunc) http.Handler {
	return middlewares.AuthMiddleware(handlerFunc)
}

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"task2-go-microservice/handlers"
	"task2-go-microservice/metrics"
	"task2-go-microservice/utils"

	"github.com/gorilla/mux"
)

func main() {
	logger := utils.NewLogger()
	userService := handlers.NewUserHandler(logger)
	integrationHandler := handlers.NewIntegrationHandler(logger)

	router := mux.NewRouter()
	router.Use(utils.RateLimitMiddleware)
	router.Use(metrics.MetricsMiddleware)

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/metrics", metrics.Handler()).Methods(http.MethodGet)

	router.HandleFunc("/api/users", userService.ListUsers).Methods(http.MethodGet)
	router.HandleFunc("/api/users/{id}", userService.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/api/users", userService.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/users/{id}", userService.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/api/users/{id}", userService.DeleteUser).Methods(http.MethodDelete)

	router.HandleFunc("/api/integration/minio/upload", integrationHandler.UploadSample).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:           router,
		Addr:              ":8080",
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	_ = os.Getenv("PLACEHOLDER")
}

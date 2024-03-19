package app

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Niiazgulov/vkquest/internal/config"
	"github.com/Niiazgulov/vkquest/internal/handlers"
	"github.com/Niiazgulov/vkquest/internal/storage"
)

func Start() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	config.Cfg = *cfg
	repo, err := storage.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/", func(r chi.Router) {
		r.Post("/adduser", handlers.CreateUserHandler(repo))
		r.Post("/addquest", handlers.CreateQuestHandler(repo))
		r.Post("/action", handlers.NewActionHandler(repo))
		r.Get("/{id}", handlers.GetUserInfoHandler(repo))
	})

	http.ListenAndServe(config.Cfg.ServerAddress, r)
}

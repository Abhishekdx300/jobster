package api

import (
	"net/http"

	"github.com/Abhishekdx300/jobster/internal/handlers"
	"github.com/Abhishekdx300/jobster/internal/repositories"
	"github.com/Abhishekdx300/jobster/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(db *mongo.Database) http.Handler {
	r := chi.NewRouter()

	jobRepo := repositories.NewJobRepository(db)
	applyRepo := repositories.NewApplyRepository(db)
	peopleRepo := repositories.NewPeopleRepository(db)

	jobService := services.NewJobService(jobRepo, applyRepo, peopleRepo)
	applyService := services.NewApplyService(applyRepo, jobRepo)
	peopleService := services.NewPeopleService(peopleRepo, jobRepo)

	jobHandler := handlers.NewJobHandler(jobService)
	applyHandler := handlers.NewApplyHandler(applyService)
	peopleHandler := handlers.NewPeopleHandler(peopleService)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/jobs", func(r chi.Router) {
		r.Post("/", jobHandler.Create)
		r.Get("/search", jobHandler.Search)

		r.Route("/{jobId}", func(r chi.Router) {
			r.Get("/", jobHandler.GetById)
			r.Put("/", jobHandler.Update)
			r.Delete("/", jobHandler.Delete)

			// Routes for Applies
			r.Route("/applies", func(r chi.Router) {
				r.Get("/", applyHandler.FindByJobId)
				r.Post("/", applyHandler.Create)
				r.Put("/{applyId}", applyHandler.Update)
				r.Delete("/{applyId}", applyHandler.Delete)
			})

			// Routes for People Reached
			r.Route("/people", func(r chi.Router) {
				r.Get("/", peopleHandler.FindByJobId)
				r.Post("/", peopleHandler.Create)
				r.Put("/{peopleId}", peopleHandler.Update)
				r.Delete("/{peopleId}", peopleHandler.Delete)
			})
		})

	})

	return r
}

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/isobekeito/safearea/internal/controller"
	"github.com/isobekeito/safearea/internal/infrastructure/external"
	"github.com/isobekeito/safearea/internal/infrastructure/persistence"
	"github.com/isobekeito/safearea/internal/usecase"
)

var chiLambda *chiadapter.ChiLambda

// Controllers
var (
	searchController    *controller.SearchController
	disastersController *controller.DisastersController
)

func init() {
	// Logger setup
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database setup
	db, err := persistence.InitDB()
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
	}

	// DI setup
	if db != nil {
		geocodingClient := external.NewGeocodingClient()
		disasterRepo := persistence.NewDisasterRepository(db)

		searchUseCase := usecase.NewSearchUseCase(geocodingClient, disasterRepo)
		disasterHistoryUseCase := usecase.NewDisasterHistoryUseCase(disasterRepo)

		searchController = controller.NewSearchController(searchUseCase)
		disastersController = controller.NewDisastersController(disasterHistoryUseCase)
	}

	// Router setup
	r := newRouter()
	chiLambda = chiadapter.New(r)
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: 本番環境では適切に設定
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", controller.HealthCheck)

		if searchController != nil {
			r.Get("/search", searchController.Search)
		}
		if disastersController != nil {
			r.Get("/search/disasters", disastersController.GetDisasters)
		}
	})

	return r
}

// Lambda handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Check if running in Lambda
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(Handler)
	} else {
		// Local development
		slog.Info("Starting local server on :8080")
		r := newRouter()
		http.ListenAndServe(":8080", r)
	}
}

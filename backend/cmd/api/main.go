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
)

var chiLambda *chiadapter.ChiLambda

var routers = map[string]http.HandlerFunc{
	"/health":    controller.HealthCheck,
	"/risk":      controller.GetRisk,
	"/disasters": controller.GetDisasters,
}

func init() {
	// Logger setup
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

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
	Route(r)
	return r
}

// Lambda handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func Route(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		for key, value := range routers {
			r.Get(key, value)
		}
	})
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

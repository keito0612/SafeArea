module github.com/isobekeito/safearea

go 1.22

require (
	// Lambda
	github.com/aws/aws-lambda-go v1.47.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.16.2

	// Router
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-chi/cors v1.2.1

	// ORM
	gorm.io/gorm v1.25.12
	gorm.io/driver/postgres v1.5.9

	// Validation
	github.com/go-playground/validator/v10 v10.22.1

	// UUID
	github.com/google/uuid v1.6.0

	// Environment (local development)
	github.com/joho/godotenv v1.5.1
)

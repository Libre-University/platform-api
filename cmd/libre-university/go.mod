module github.com/Libre-University/platform-api/cmd/libre-university

go 1.26.0

replace (
	github.com/Libre-University/platform-api/core => ../../core
	github.com/Libre-University/platform-api/modules/lms => ../../modules/lms
	github.com/Libre-University/platform-api/modules/obs => ../../modules/obs
	github.com/Libre-University/platform-api/platform => ../../platform
)

require (
	github.com/Libre-University/platform-api/core v0.0.0-00010101000000-000000000000
	github.com/Libre-University/platform-api/modules/lms v0.0.0-00010101000000-000000000000
	github.com/Libre-University/platform-api/modules/obs v0.0.0-00010101000000-000000000000
	github.com/Libre-University/platform-api/platform v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-chi/chi/v5 v5.3.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.11.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pressly/goose/v3 v3.28.0 // indirect
	github.com/sethvargo/go-retry v0.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/text v0.41.0 // indirect
)

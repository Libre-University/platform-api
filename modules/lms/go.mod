module github.com/Libre-University/platform-api/modules/lms

go 1.26.0

replace github.com/Libre-University/platform-api/core => ../../core

require (
	github.com/Libre-University/platform-api/core v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.3.2
)

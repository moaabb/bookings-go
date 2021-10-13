package dbrepo

import (
	"database/sql"

	"github.com/moaabb/bookings-go/internal/config"
	"github.com/moaabb/bookings-go/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testingDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(d *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  d,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testingDBRepo{
		App: a,
	}
}

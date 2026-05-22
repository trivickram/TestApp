package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var errAlreadyLinked = errors.New("doctor already linked to this clinic")

type store struct {
	db *sql.DB
}

type clinic struct {
	id   int64
	name string
}

type doctor struct {
	id             int64
	name           string
	specialization string
}

type patient struct {
	id   int64
	name string
	age  int32
}

type appointment struct {
	id          int64
	clinicID    int64
	doctorID    int64
	patientID   int64
	scheduledAt string
	status      string
}

func newStore(dsn string) (*store, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &store{db: db}, nil
}

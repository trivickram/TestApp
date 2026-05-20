package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type store struct {
	db *sql.DB
}

func newStore(dsn string) (*store, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS patients (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			age INT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS appointments (
			id VARCHAR(36) PRIMARY KEY,
			patient_id VARCHAR(36) NOT NULL,
			doctor VARCHAR(255) NOT NULL,
			scheduled_at VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL,
			INDEX idx_patient_id (patient_id),
			INDEX idx_status (status)
		)
	`)
	if err != nil {
		return nil, err
	}
	return &store{db: db}, nil
}

func (s *store) insertPatient(id, name string, age int32) error {
	_, err := s.db.Exec("INSERT INTO patients (id, name, age) VALUES (?, ?, ?)", id, name, age)
	return err
}

func (s *store) patientExists(id string) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM patients WHERE id = ?", id).Scan(&count)
	return count > 0, err
}

func (s *store) insertAppointment(id, patientID, doctor, scheduledAt, status string) error {
	_, err := s.db.Exec(
		"INSERT INTO appointments (id, patient_id, doctor, scheduled_at, status) VALUES (?, ?, ?, ?, ?)",
		id, patientID, doctor, scheduledAt, status,
	)
	return err
}

func (s *store) getAppointmentStatus(id string) (string, bool, error) {
	var st string
	err := s.db.QueryRow("SELECT status FROM appointments WHERE id = ?", id).Scan(&st)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	return st, err == nil, err
}

func (s *store) updateAppointmentStatus(id, status string) (bool, error) {
	res, err := s.db.Exec("UPDATE appointments SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

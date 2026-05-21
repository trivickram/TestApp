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

func (s *store) listClinics() ([]clinic, error) {
	rows, err := s.db.Query("SELECT id, name FROM clinics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []clinic
	for rows.Next() {
		var c clinic
		if err := rows.Scan(&c.id, &c.name); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *store) insertDoctor(name, spec string) (int64, error) {
	res, err := s.db.Exec("INSERT INTO doctors (name, specialization) VALUES (?, ?)", name, spec)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *store) linkDoctor(clinicID, doctorID int64) error {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM clinic_doctors WHERE clinic_id = ? AND doctor_id = ?", clinicID, doctorID).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return errAlreadyLinked
	}
	_, err := s.db.Exec("INSERT INTO clinic_doctors (clinic_id, doctor_id) VALUES (?, ?)", clinicID, doctorID)
	return err
}

func (s *store) searchDoctors(q string, clinicID int64) ([]doctor, error) {
	var rows *sql.Rows
	var err error
	if clinicID > 0 {
		sql := `SELECT d.id, d.name, d.specialization FROM doctors d
			JOIN clinic_doctors cd ON cd.doctor_id = d.id
			WHERE cd.clinic_id = ? AND d.name LIKE ? LIMIT 20`
		rows, err = s.db.Query(sql, clinicID, "%"+q+"%")
	} else {
		rows, err = s.db.Query("SELECT id, name, specialization FROM doctors WHERE name LIKE ? LIMIT 20", "%"+q+"%")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []doctor
	for rows.Next() {
		var d doctor
		if err := rows.Scan(&d.id, &d.name, &d.specialization); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *store) searchPatients(q string) ([]patient, error) {
	rows, err := s.db.Query("SELECT id, name, age FROM patients WHERE name LIKE ? LIMIT 20", "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []patient
	for rows.Next() {
		var p patient
		if err := rows.Scan(&p.id, &p.name, &p.age); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *store) listClinicDoctors(clinicID int64) ([]doctor, error) {
	rows, err := s.db.Query(`
		SELECT d.id, d.name, d.specialization
		FROM doctors d
		JOIN clinic_doctors cd ON cd.doctor_id = d.id
		WHERE cd.clinic_id = ?
	`, clinicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []doctor
	for rows.Next() {
		var d doctor
		if err := rows.Scan(&d.id, &d.name, &d.specialization); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *store) insertPatient(name string, age int32) (int64, error) {
	res, err := s.db.Exec("INSERT INTO patients (name, age) VALUES (?, ?)", name, age)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *store) listPatients() ([]patient, error) {
	rows, err := s.db.Query("SELECT id, name, age FROM patients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []patient
	for rows.Next() {
		var p patient
		if err := rows.Scan(&p.id, &p.name, &p.age); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *store) doctorConflict(doctorID int64, scheduledAt string) (bool, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM appointments WHERE doctor_id = ? AND scheduled_at = ? AND status = 'SCHEDULED'",
		doctorID, scheduledAt,
	).Scan(&count)
	return count > 0, err
}

func (s *store) patientConflict(patientID int64, scheduledAt string) (bool, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM appointments WHERE patient_id = ? AND scheduled_at = ? AND status = 'SCHEDULED'",
		patientID, scheduledAt,
	).Scan(&count)
	return count > 0, err
}

func (s *store) updateAppointmentStatus(id int64, newStatus string) (*appointment, error) {
	res, err := s.db.Exec("UPDATE appointments SET status = ? WHERE id = ?", newStatus, id)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, sql.ErrNoRows
	}
	var a appointment
	err = s.db.QueryRow(
		"SELECT id, clinic_id, doctor_id, patient_id, DATE_FORMAT(scheduled_at, '%Y-%m-%dT%H:%i:%s'), status FROM appointments WHERE id = ?",
		id,
	).Scan(&a.id, &a.clinicID, &a.doctorID, &a.patientID, &a.scheduledAt, &a.status)
	return &a, err
}

func (s *store) insertAppointment(clinicID, doctorID, patientID int64, scheduledAt string) (int64, error) {
	res, err := s.db.Exec(
		"INSERT INTO appointments (clinic_id, doctor_id, patient_id, scheduled_at, status) VALUES (?, ?, ?, ?, 'SCHEDULED')",
		clinicID, doctorID, patientID, scheduledAt,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *store) listAppointments(clinicID, doctorID, patientID int64, date string) ([]appointment, error) {
	q := `SELECT id, clinic_id, doctor_id, patient_id, DATE_FORMAT(scheduled_at, '%Y-%m-%dT%H:%i:%s'), status FROM appointments WHERE 1=1`
	args := []any{}
	if clinicID > 0 {
		q += " AND clinic_id = ?"
		args = append(args, clinicID)
	}
	if doctorID > 0 {
		q += " AND doctor_id = ?"
		args = append(args, doctorID)
	}
	if patientID > 0 {
		q += " AND patient_id = ?"
		args = append(args, patientID)
	}
	if date != "" {
		q += " AND DATE(scheduled_at) = ?"
		args = append(args, date)
	}
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []appointment
	for rows.Next() {
		var a appointment
		if err := rows.Scan(&a.id, &a.clinicID, &a.doctorID, &a.patientID, &a.scheduledAt, &a.status); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

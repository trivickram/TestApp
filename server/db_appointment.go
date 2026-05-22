package main

import "database/sql"

func (s *store) doctorHasPendingAppointment(doctorID int64) (bool, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM appointments WHERE doctor_id = ? AND status = 'SCHEDULED'",
		doctorID,
	).Scan(&count)
	return count > 0, err
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

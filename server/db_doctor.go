package main

import "database/sql"

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

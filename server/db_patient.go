package main

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

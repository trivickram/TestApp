package main

func (s *store) listClinics() ([]clinic, error) {
	rows, err := s.db.Query(`SELECT id, name FROM clinics ORDER BY id`)
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

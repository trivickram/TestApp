USE hospital;

EXPLAIN SELECT * FROM patients WHERE id = 'p-001';

EXPLAIN ANALYZE SELECT * FROM patients WHERE id = 'p-001';


EXPLAIN SELECT * FROM patients WHERE age BETWEEN 30 AND 50;

EXPLAIN ANALYZE SELECT * FROM patients WHERE age BETWEEN 30 AND 50;

EXPLAIN SELECT * FROM appointments WHERE patient_id = 'p-001';

EXPLAIN ANALYZE SELECT * FROM appointments WHERE patient_id = 'p-001';


EXPLAIN SELECT * FROM appointments WHERE status = 'SCHEDULED';

EXPLAIN ANALYZE SELECT * FROM appointments WHERE status = 'SCHEDULED';

EXPLAIN SELECT * FROM appointments WHERE doctor = 'Dr. Smith';

EXPLAIN ANALYZE SELECT * FROM appointments WHERE doctor = 'Dr. Smith';

EXPLAIN SELECT p.name, p.age, a.doctor, a.status
FROM patients p
JOIN appointments a ON a.patient_id = p.id
WHERE a.status = 'SCHEDULED';

EXPLAIN ANALYZE SELECT p.name, p.age, a.doctor, a.status
FROM patients p
JOIN appointments a ON a.patient_id = p.id
WHERE a.status = 'SCHEDULED';


EXPLAIN SELECT p.name, a.doctor, a.status
FROM patients p
JOIN appointments a ON a.patient_id = p.id
WHERE p.age > 60;

EXPLAIN ANALYZE SELECT p.name, a.doctor, a.status
FROM patients p
JOIN appointments a ON a.patient_id = p.id
WHERE p.age > 60;


EXPLAIN SELECT status, COUNT(*) AS total
FROM appointments
GROUP BY status;

EXPLAIN ANALYZE SELECT status, COUNT(*) AS total
FROM appointments
GROUP BY status;


EXPLAIN SELECT doctor, COUNT(*) AS total
FROM appointments
GROUP BY doctor
ORDER BY total DESC;

EXPLAIN ANALYZE SELECT doctor, COUNT(*) AS total
FROM appointments
GROUP BY doctor
ORDER BY total DESC;

EXPLAIN SELECT p.name, p.age, sub.cnt
FROM patients p
JOIN (
    SELECT patient_id, COUNT(*) AS cnt
    FROM appointments
    GROUP BY patient_id
) sub ON sub.patient_id = p.id
ORDER BY sub.cnt DESC
LIMIT 10;

EXPLAIN ANALYZE SELECT p.name, p.age, sub.cnt
FROM patients p
JOIN (
    SELECT patient_id, COUNT(*) AS cnt
    FROM appointments
    GROUP BY patient_id
) sub ON sub.patient_id = p.id
ORDER BY sub.cnt DESC
LIMIT 10;


EXPLAIN FORMAT=JSON SELECT p.name, a.doctor, a.status
FROM patients p
JOIN appointments a ON a.patient_id = p.id
WHERE a.status = 'COMPLETED' AND p.age < 40;

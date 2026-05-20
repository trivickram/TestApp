USE hospital;

-- ─────────────────────────────────────────────
-- 1. PRIMARY KEY lookup  →  type: const (best)
-- ─────────────────────────────────────────────
EXPLAIN SELECT * FROM patients WHERE id = 'p-001';

EXPLAIN ANALYZE SELECT * FROM patients WHERE id = 'p-001';

-- ─────────────────────────────────────────────
-- 2. INDEX range scan on age
-- ─────────────────────────────────────────────
EXPLAIN SELECT * FROM patients WHERE age BETWEEN 30 AND 50;

EXPLAIN ANALYZE SELECT * FROM patients WHERE age BETWEEN 30 AND 50;

-- ─────────────────────────────────────────────
-- 3. INDEX lookup on appointments.patient_id
-- ─────────────────────────────────────────────
EXPLAIN SELECT * FROM appointments WHERE patient_id = 'p-001';

EXPLAIN ANALYZE SELECT * FROM appointments WHERE patient_id = 'p-001';

-- ─────────────────────────────────────────────
-- 4. INDEX lookup on appointments.status
-- ─────────────────────────────────────────────
EXPLAIN SELECT * FROM appointments WHERE status = 'SCHEDULED';

EXPLAIN ANALYZE SELECT * FROM appointments WHERE status = 'SCHEDULED';

-- ─────────────────────────────────────────────
-- 5. FULL TABLE SCAN (no index on doctor)
--    compare cost vs indexed queries above
-- ─────────────────────────────────────────────
EXPLAIN SELECT * FROM appointments WHERE doctor = 'Dr. Smith';

EXPLAIN ANALYZE SELECT * FROM appointments WHERE doctor = 'Dr. Smith';

-- ─────────────────────────────────────────────
-- 6. JOIN  →  nested loop, index on both sides
-- ─────────────────────────────────────────────
EXPLAIN SELECT p.name, p.age, a.doctor, a.status
FROM patients p
JOIN appointments a ON a.patient_id = p.id
WHERE a.status = 'SCHEDULED';

EXPLAIN ANALYZE SELECT p.name, p.age, a.doctor, a.status
FROM patients p
JOIN appointments a ON a.patient_id = p.id
WHERE a.status = 'SCHEDULED';

-- ─────────────────────────────────────────────
-- 7. JOIN with age range filter
-- ─────────────────────────────────────────────
EXPLAIN SELECT p.name, a.doctor, a.status
FROM patients p
JOIN appointments a ON a.patient_id = p.id
WHERE p.age > 60;

EXPLAIN ANALYZE SELECT p.name, a.doctor, a.status
FROM patients p
JOIN appointments a ON a.patient_id = p.id
WHERE p.age > 60;

-- ─────────────────────────────────────────────
-- 8. Aggregation  →  COUNT per status
-- ─────────────────────────────────────────────
EXPLAIN SELECT status, COUNT(*) AS total
FROM appointments
GROUP BY status;

EXPLAIN ANALYZE SELECT status, COUNT(*) AS total
FROM appointments
GROUP BY status;

-- ─────────────────────────────────────────────
-- 9. Aggregation  →  appointments per doctor
-- ─────────────────────────────────────────────
EXPLAIN SELECT doctor, COUNT(*) AS total
FROM appointments
GROUP BY doctor
ORDER BY total DESC;

EXPLAIN ANALYZE SELECT doctor, COUNT(*) AS total
FROM appointments
GROUP BY doctor
ORDER BY total DESC;

-- ─────────────────────────────────────────────
-- 10. Subquery  →  patients with most appointments
-- ─────────────────────────────────────────────
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

-- ─────────────────────────────────────────────
-- 11. JSON format  →  full optimizer detail
-- ─────────────────────────────────────────────
EXPLAIN FORMAT=JSON SELECT p.name, a.doctor, a.status
FROM patients p
JOIN appointments a ON a.patient_id = p.id
WHERE a.status = 'COMPLETED' AND p.age < 40;

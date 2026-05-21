-- =============================================================
-- performance.sql  ·  Query performance analysis
-- Requires: MySQL 8.0.18+  (EXPLAIN ANALYZE)
-- Run after seed.sql:
--   mysql -u root -pPassword@123 hospital < scripts/performance.sql
-- =============================================================

USE hospital;

SELECT '=======================================================' AS '';
SELECT 'PERFORMANCE ANALYSIS — Hospital DB'                       AS '';
SELECT '=======================================================' AS '';

-- -------------------------------------------------------------
-- Q1. INDEX SCAN: appointments for a doctor in a date range
--     Expected: uses idx_doctor_time  (doctor_id, scheduled_at)
-- -------------------------------------------------------------
SELECT '--- Q1: Doctor appointments in date range ---' AS query;
EXPLAIN ANALYZE
SELECT id, patient_id, scheduled_at, status
FROM   appointments
WHERE  doctor_id    = 5
  AND  scheduled_at BETWEEN '2025-01-01' AND '2025-06-30';

-- -------------------------------------------------------------
-- Q2. INDEX SCAN: appointments for a patient in a date range
--     Expected: uses idx_patient_time  (patient_id, scheduled_at)
-- -------------------------------------------------------------
SELECT '--- Q2: Patient appointments in date range ---' AS query;
EXPLAIN ANALYZE
SELECT id, doctor_id, clinic_id, scheduled_at, status
FROM   appointments
WHERE  patient_id   = 10
  AND  scheduled_at BETWEEN '2025-01-01' AND '2025-12-31';

-- -------------------------------------------------------------
-- Q3. COMPOSITE PK LOOKUP: doctors in a clinic
--     Expected: uses PRIMARY KEY on clinic_doctors (clinic_id, doctor_id)
-- -------------------------------------------------------------
SELECT '--- Q3: List doctors for a clinic ---' AS query;
EXPLAIN ANALYZE
SELECT d.id, d.name, d.specialization
FROM   doctors d
JOIN   clinic_doctors cd ON cd.doctor_id = d.id
WHERE  cd.clinic_id = 2;

-- -------------------------------------------------------------
-- Q4. AGGREGATE + JOIN: appointment count per doctor per clinic
--     Expected: index on clinic_id for appointments (none yet → watch rows)
-- -------------------------------------------------------------
SELECT '--- Q4: Appointment count per doctor for clinic 1 ---' AS query;
EXPLAIN ANALYZE
SELECT d.name, d.specialization, COUNT(*) AS appt_count
FROM   appointments a
JOIN   doctors d ON d.id = a.doctor_id
WHERE  a.clinic_id = 1
GROUP  BY d.id, d.name, d.specialization
ORDER  BY appt_count DESC;

-- -------------------------------------------------------------
-- Q5. FULL TABLE SCAN CANDIDATE: filter by status only
--     Expected: full scan (no index on `status`) → candidate for new index
-- -------------------------------------------------------------
SELECT '--- Q5: Count appointments by status (no index on status) ---' AS query;
EXPLAIN ANALYZE
SELECT status, COUNT(*) AS cnt
FROM   appointments
GROUP  BY status;

-- -------------------------------------------------------------
-- Q6. FULL TABLE SCAN CANDIDATE: doctor name LIKE with leading wildcard
--     Expected: full scan — leading % prevents index use
-- -------------------------------------------------------------
SELECT '--- Q6: Doctor name search with leading wildcard (full scan) ---' AS query;
EXPLAIN ANALYZE
SELECT id, name, specialization
FROM   doctors
WHERE  name LIKE '%Cardiology%';

-- -------------------------------------------------------------
-- Q7. PARTIAL INDEX USE: doctor name LIKE with trailing wildcard only
--     Expected: range scan if there is an index on name; otherwise full scan
-- -------------------------------------------------------------
SELECT '--- Q7: Doctor name prefix search (trailing wildcard) ---' AS query;
EXPLAIN ANALYZE
SELECT id, name, specialization
FROM   doctors
WHERE  name LIKE 'Dr. A%';

-- -------------------------------------------------------------
-- Q8. AGGREGATE: busiest clinics overall
-- -------------------------------------------------------------
SELECT '--- Q8: Total appointments per clinic ---' AS query;
EXPLAIN ANALYZE
SELECT c.name, COUNT(a.id) AS total
FROM   clinics c
LEFT   JOIN appointments a ON a.clinic_id = c.id
GROUP  BY c.id, c.name
ORDER  BY total DESC;

-- -------------------------------------------------------------
-- Q9. DATE FUNCTION: appointments on a specific calendar date
--     Expected: may prevent index use if optimizer can't use DATE()
-- -------------------------------------------------------------
SELECT '--- Q9: Appointments on a specific date (DATE function) ---' AS query;
EXPLAIN ANALYZE
SELECT id, doctor_id, patient_id, status
FROM   appointments
WHERE  DATE(scheduled_at) = '2025-03-15';

-- Q9b — equivalent rewrite using a range (optimizer-friendly)
SELECT '--- Q9b: Same query rewritten as range (index-friendly) ---' AS query;
EXPLAIN ANALYZE
SELECT id, doctor_id, patient_id, status
FROM   appointments
WHERE  scheduled_at >= '2025-03-15 00:00:00'
  AND  scheduled_at <  '2025-03-16 00:00:00';

-- -------------------------------------------------------------
-- Q10. CORRELATED: patients who have more than 1 appointment
-- -------------------------------------------------------------
SELECT '--- Q10: Patients with more than 1 appointment ---' AS query;
EXPLAIN ANALYZE
SELECT p.id, p.name, COUNT(a.id) AS appts
FROM   patients p
JOIN   appointments a ON a.patient_id = p.id
GROUP  BY p.id, p.name
HAVING COUNT(a.id) > 1
ORDER  BY appts DESC
LIMIT  20;

-- =============================================================
-- Index inventory
-- =============================================================
SELECT '=======================================================' AS '';
SELECT 'CURRENT INDEXES'                                          AS '';
SELECT '=======================================================' AS '';

SHOW INDEX FROM appointments;
SHOW INDEX FROM doctors;
SHOW INDEX FROM patients;
SHOW INDEX FROM clinic_doctors;

-- =============================================================
-- Suggested indexes based on the queries above
-- =============================================================
SELECT '=======================================================' AS '';
SELECT 'SUGGESTED ADDITIONAL INDEXES'                             AS '';
SELECT '=======================================================' AS '';

SELECT
  'ALTER TABLE appointments ADD INDEX idx_clinic_id (clinic_id);'
    AS suggestion,
  'Speeds up Q4/Q8 — clinic-scoped appointment queries' AS reason
UNION ALL SELECT
  'ALTER TABLE appointments ADD INDEX idx_status (status);',
  'Speeds up Q5 — status-based filtering/grouping'
UNION ALL SELECT
  'ALTER TABLE doctors ADD INDEX idx_doctor_name (name);',
  'Speeds up Q7 — prefix-based name searches'
UNION ALL SELECT
  'ALTER TABLE patients ADD INDEX idx_patient_name (name);',
  'Speeds up patient name prefix searches';

-- =============================================================
-- seed.sql  ·  High-volume population script
-- Requires: MySQL 8.0+  ·  schema already created via init.sql
-- Run: mysql -u root -pPassword@123 hospital < scripts/seed.sql
-- Inserts:
--   10  clinics  (3 original + 7 new)
-- 1000  doctors  (100 per clinic)
-- 5000  patients
-- 1000  clinic-doctor links (1 clinic per doctor)
-- 20000  appointments (all COMPLETED — doctors are free to schedule)
-- =============================================================

USE hospital;

SET SESSION cte_max_recursion_depth = 20001;

-- ------------------------------------------------------------
-- Wipe existing seed data
-- ------------------------------------------------------------
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE appointments;
TRUNCATE TABLE clinic_doctors;
TRUNCATE TABLE patients;
TRUNCATE TABLE doctors;
DELETE FROM clinics WHERE id > 3;
SET FOREIGN_KEY_CHECKS = 1;

ALTER TABLE clinics      AUTO_INCREMENT = 4;
ALTER TABLE doctors      AUTO_INCREMENT = 1;
ALTER TABLE patients     AUTO_INCREMENT = 1;
ALTER TABLE appointments AUTO_INCREMENT = 1;

-- ------------------------------------------------------------
-- 1. Clinics  (7 new → total 10)
-- ------------------------------------------------------------
INSERT INTO clinics (name) VALUES
  ('Sunrise Health Centre'),
  ('Green Valley Hospital'),
  ('City Medical Institute'),
  ('Lakeview Clinic'),
  ('Metro General Hospital'),
  ('Heritage Medical Centre'),
  ('Pearl Diagnostics & Care');

-- ------------------------------------------------------------
-- 2. Doctors  (1000 records, 10 specializations cycling)
-- ------------------------------------------------------------
INSERT INTO doctors (name, specialization)
WITH RECURSIVE seq(n) AS (
  SELECT 1 UNION ALL SELECT n + 1 FROM seq WHERE n < 1000
)
SELECT
  CONCAT('Dr. ',
    CHAR(65 + ((n - 1) MOD 26)),
    CHAR(65 + ((n - 1) DIV 26 MOD 26)),
    LPAD(n, 4, '0')),
  ELT(1 + ((n - 1) MOD 10),
    'Cardiology','Neurology','Orthopedics','Pediatrics','Dermatology',
    'Oncology','Gynecology','Psychiatry','Radiology','Urology')
FROM seq;

-- ------------------------------------------------------------
-- 3. Patients  (5000 records, ages 18–79)
-- ------------------------------------------------------------
INSERT INTO patients (name, age)
WITH RECURSIVE seq(n) AS (
  SELECT 1 UNION ALL SELECT n + 1 FROM seq WHERE n < 5000
)
SELECT
  CONCAT('Patient ', LPAD(n, 4, '0')),
  18 + ((n - 1) MOD 62)
FROM seq;

-- ------------------------------------------------------------
-- 4. Clinic-Doctor links  (1 clinic per doctor, 100 doctors per clinic)
--    doctor  1– 100 → clinic 1
--    doctor 101– 200 → clinic 2  … doctor 901–1000 → clinic 10
-- ------------------------------------------------------------
INSERT INTO clinic_doctors (clinic_id, doctor_id)
WITH RECURSIVE seq(n) AS (
  SELECT 1 UNION ALL SELECT n + 1 FROM seq WHERE n < 1000
)
SELECT ((n - 1) DIV 100) + 1, n FROM seq;

-- ------------------------------------------------------------
-- 5. Appointments  (20 000 records, all COMPLETED)
--
--    Row n (0-based):
--      doctor_id    = (n MOD 1000) + 1      → each doctor gets 20 appts
--      patient_id   = (n MOD 5000) + 1      → each patient gets 4 appts
--      clinic_id    = ((n MOD 1000) DIV 100) + 1
--      scheduled_at = 2020-01-01 00:00 + n hours
--
--    No doctor conflict: same doctor reappears every 1000 rows
--      → 1000 hours apart → always different datetime ✓
--    No patient conflict: same patient reappears every 5000 rows
--      → 5000 hours apart → always different datetime ✓
-- ------------------------------------------------------------
INSERT INTO appointments (clinic_id, doctor_id, patient_id, scheduled_at, status)
WITH RECURSIVE seq(n) AS (
  SELECT 0 UNION ALL SELECT n + 1 FROM seq WHERE n < 19999
)
SELECT
  ((n MOD 1000) DIV 100) + 1,
  (n MOD 1000) + 1,
  (n MOD 5000) + 1,
  DATE_ADD(CAST('2020-01-01 00:00:00' AS DATETIME), INTERVAL n HOUR),
  'COMPLETED'
FROM seq;

-- ------------------------------------------------------------
-- 6. SCHEDULED appointments (1 per clinic, future date)
--    First doctor of each clinic batch: 1, 101, 201, ..., 901
--    Each doctor's last appointment is COMPLETED above ✓
-- ------------------------------------------------------------
INSERT INTO appointments (clinic_id, doctor_id, patient_id, scheduled_at, status) VALUES
  (1,   1,    1,  '2026-06-10 09:00:00', 'SCHEDULED'),
  (2,   101,  2,  '2026-06-10 09:00:00', 'SCHEDULED'),
  (3,   201,  3,  '2026-06-10 09:00:00', 'SCHEDULED'),
  (4,   301,  4,  '2026-06-10 09:00:00', 'SCHEDULED'),
  (5,   401,  5,  '2026-06-10 09:00:00', 'SCHEDULED'),
  (6,   501,  6,  '2026-06-10 09:00:00', 'SCHEDULED'),
  (7,   601,  7,  '2026-06-10 09:00:00', 'SCHEDULED'),
  (8,   701,  8,  '2026-06-10 09:00:00', 'SCHEDULED'),
  (9,   801,  9,  '2026-06-10 09:00:00', 'SCHEDULED'),
  (10,  901, 10,  '2026-06-10 09:00:00', 'SCHEDULED');

-- ------------------------------------------------------------
-- Summary
-- ------------------------------------------------------------
SELECT 'clinics'        AS `table`, COUNT(*) AS `rows` FROM clinics
UNION ALL SELECT 'doctors',        COUNT(*) FROM doctors
UNION ALL SELECT 'patients',       COUNT(*) FROM patients
UNION ALL SELECT 'clinic_doctors', COUNT(*) FROM clinic_doctors
UNION ALL SELECT 'appointments',   COUNT(*) FROM appointments;


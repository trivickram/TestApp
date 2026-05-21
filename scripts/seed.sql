-- =============================================================
-- seed.sql  ·  Mass population script
-- Requires: MySQL 8.0+  ·  schema already created via init.sql
-- Run: mysql -u root -pPassword@123 hospital < scripts/seed.sql
-- Inserts:
--   10  clinics  (7 new + 3 from init.sql)
--  200  doctors
-- 5000  patients
--  ~400  clinic-doctor links  (each doctor → 2 clinics)
-- 10000  appointments         (no scheduling conflicts)
-- =============================================================

USE hospital;

-- Allow deep recursion for the number-generator CTEs
SET SESSION cte_max_recursion_depth = 10001;

-- ------------------------------------------------------------
-- Wipe existing seed data; preserve schema + original 3 clinics
-- ------------------------------------------------------------
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE appointments;
TRUNCATE TABLE clinic_doctors;
TRUNCATE TABLE patients;
TRUNCATE TABLE doctors;
DELETE FROM clinics WHERE id > 3;
SET FOREIGN_KEY_CHECKS = 1;

-- Reset auto-increment counters
ALTER TABLE clinics      AUTO_INCREMENT = 4;
ALTER TABLE doctors      AUTO_INCREMENT = 1;
ALTER TABLE patients     AUTO_INCREMENT = 1;
ALTER TABLE appointments AUTO_INCREMENT = 1;

-- ------------------------------------------------------------
-- 1. Clinics  (7 more → total 10)
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
-- 2. Doctors  (200 records, 10 specializations cycling)
-- ------------------------------------------------------------
INSERT INTO doctors (name, specialization)
WITH RECURSIVE seq(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 200
)
SELECT
  CONCAT(
    'Dr. ',
    CHAR(64 + ((n - 1) MOD 26) + 1),   -- letter A-Z
    CHAR(64 + ((n - 1) DIV 26 MOD 26) + 1),
    LPAD(n, 3, '0')
  ),
  ELT(
    1 + ((n - 1) MOD 10),
    'Cardiology', 'Neurology', 'Orthopedics', 'Pediatrics', 'Dermatology',
    'Oncology',   'Gynecology', 'Psychiatry', 'Radiology',  'Urology'
  )
FROM seq;

-- ------------------------------------------------------------
-- 3. Patients  (5 000 records, ages 18-79)
-- ------------------------------------------------------------
INSERT INTO patients (name, age)
WITH RECURSIVE seq(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 5000
)
SELECT
  CONCAT('Patient ', LPAD(n, 4, '0')),
  18 + ((n - 1) MOD 62)   -- 18..79
FROM seq;

-- ------------------------------------------------------------
-- 4. Clinic-Doctor links  (each doctor linked to 2 distinct clinics)
--    Clinic A = ((doctor_id - 1) MOD 10) + 1
--    Clinic B = ((doctor_id)     MOD 10) + 1   (next clinic, wraps)
-- ------------------------------------------------------------
INSERT IGNORE INTO clinic_doctors (clinic_id, doctor_id)
WITH RECURSIVE seq(n) AS (
  SELECT 1
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 200
)
SELECT ((n - 1) MOD 10) + 1, n FROM seq
UNION ALL
SELECT ( n      MOD 10) + 1, n FROM seq;

-- ------------------------------------------------------------
-- 5. Appointments  (10 000 records, zero scheduling conflicts)
--
--    Row n (0-based) gets:
--      doctor_id    = (n MOD 200)  + 1   → each doctor gets 50 appts
--      patient_id   = (n MOD 5000) + 1   → each patient gets 2 appts
--      clinic_id    = (n MOD 10)   + 1
--      scheduled_at = 2025-01-01 + floor(n/24) days + (n MOD 24) hours
--
--    Doctor conflict check:
--      Doctor d reappears at n = d-1, d-1+200, d-1+400, ...
--      Each step advances time by exactly 200 hours → no same-time duplicate.
--
--    Patient conflict check:
--      Patient p reappears at n = p-1, p-1+5000, ...
--      Each step advances time by 5000 hours → no same-time duplicate.
-- ------------------------------------------------------------
INSERT INTO appointments (clinic_id, doctor_id, patient_id, scheduled_at, status)
WITH RECURSIVE seq(n) AS (
  SELECT 0
  UNION ALL
  SELECT n + 1 FROM seq WHERE n < 9999
)
SELECT
  (n MOD 10)   + 1,
  (n MOD 200)  + 1,
  (n MOD 5000) + 1,
  DATE_ADD(
    DATE_ADD(CAST('2025-01-01' AS DATETIME), INTERVAL (n DIV 24) DAY),
    INTERVAL (n MOD 24) HOUR
  ),
  ELT(1 + (n MOD 3), 'SCHEDULED', 'SCHEDULED', 'COMPLETED')  -- ~33% COMPLETED
FROM seq;

-- ------------------------------------------------------------
-- Summary
-- ------------------------------------------------------------
SELECT 'clinics'        AS `table`, COUNT(*) AS `rows` FROM clinics
UNION ALL SELECT 'doctors',        COUNT(*) FROM doctors
UNION ALL SELECT 'patients',       COUNT(*) FROM patients
UNION ALL SELECT 'clinic_doctors', COUNT(*) FROM clinic_doctors
UNION ALL SELECT 'appointments',   COUNT(*) FROM appointments;

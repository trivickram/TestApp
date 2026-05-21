CREATE DATABASE IF NOT EXISTS hospital;
USE hospital;

SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS clinic_doctors;
DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS doctors;
DROP TABLE IF EXISTS clinics;
SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE clinics (
    id   INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE doctors (
    id             INT AUTO_INCREMENT PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    specialization VARCHAR(255) NOT NULL
);

CREATE TABLE patients (
    id   INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age  INT NOT NULL
);

CREATE TABLE clinic_doctors (
    clinic_id INT NOT NULL,
    doctor_id INT NOT NULL,
    PRIMARY KEY (clinic_id, doctor_id),
    FOREIGN KEY (clinic_id) REFERENCES clinics(id),
    FOREIGN KEY (doctor_id) REFERENCES doctors(id)
);

CREATE TABLE appointments (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    clinic_id    INT NOT NULL,
    doctor_id    INT NOT NULL,
    patient_id   INT NOT NULL,
    scheduled_at DATETIME NOT NULL,
    status       VARCHAR(50) NOT NULL DEFAULT 'SCHEDULED',
    FOREIGN KEY (clinic_id)  REFERENCES clinics(id),
    FOREIGN KEY (doctor_id)  REFERENCES doctors(id),
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    INDEX idx_doctor_time  (doctor_id, scheduled_at),
    INDEX idx_patient_time (patient_id, scheduled_at)
);

INSERT INTO clinics (name) VALUES ('Care and Cure'), ('Sakra Medical Hospital'), ('Apollo Hospitals');

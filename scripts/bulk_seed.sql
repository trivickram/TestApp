USE hospital;

DROP PROCEDURE IF EXISTS generate_data;

DELIMITER $$

CREATE PROCEDURE generate_data()
BEGIN
    DECLARE i       INT DEFAULT 0;
    DECLARE j       INT;
    DECLARE pid     VARCHAR(36);
    DECLARE pname   VARCHAR(255);
    DECLARE appt_count INT;

    SET SESSION bulk_insert_buffer_size = 256 * 1024 * 1024;

    WHILE i < 10000 DO
        SET pid = UUID();

        SET pname = CONCAT(
            ELT(1 + FLOOR(RAND() * 20),
                'Alice','Bob','Carol','David','Eva',
                'Frank','Grace','Henry','Isla','Jack',
                'Karen','Liam','Mia','Noah','Olivia',
                'Paul','Quinn','Rachel','Sam','Tara'),
            ' ',
            ELT(1 + FLOOR(RAND() * 20),
                'Johnson','Smith','White','Brown','Martinez',
                'Lee','Kim','Wilson','Davis','Taylor',
                'Anderson','Thomas','Jackson','Harris','Martin',
                'Garcia','Rodriguez','Lewis','Walker','Hall')
        );

        INSERT INTO patients (id, name, age)
        VALUES (pid, pname, 18 + FLOOR(RAND() * 67));

        SET appt_count = 3 + FLOOR(RAND() * 5);
        SET j = 0;

        WHILE j < appt_count DO
            INSERT INTO appointments (id, patient_id, doctor, scheduled_at, status)
            VALUES (
                UUID(),
                pid,
                ELT(1 + FLOOR(RAND() * 10),
                    'Dr. Smith','Dr. Jones','Dr. Patel','Dr. Chen','Dr. Garcia',
                    'Dr. Brown','Dr. Wilson','Dr. Davis','Dr. Martinez','Dr. Taylor'),
                DATE_FORMAT(
                    DATE_ADD('2024-01-01', INTERVAL FLOOR(RAND() * 730) DAY),
                    '%Y-%m-%dT%H:%i:%SZ'
                ),
                ELT(1 + FLOOR(RAND() * 10),
                    'SCHEDULED','SCHEDULED','SCHEDULED','SCHEDULED','SCHEDULED',
                    'SCHEDULED','COMPLETED','COMPLETED','COMPLETED','CANCELLED')
            );
            SET j = j + 1;
        END WHILE;

        SET i = i + 1;
    END WHILE;
END$$

DELIMITER ;

CALL generate_data();

DROP PROCEDURE IF EXISTS generate_data;

SELECT 'patients'    AS tbl, COUNT(*) AS total_rows FROM patients
UNION ALL
SELECT 'appointments' AS tbl, COUNT(*) AS total_rows FROM appointments;

CREATE TABLE IF NOT EXISTS late_days (
    id SERIAL PRIMARY KEY,
    student_id VARCHAR(255) NOT NULL,
    student_email VARCHAR(255) NOT NULL,
    days INT NOT NULL
);
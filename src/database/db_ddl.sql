SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS teachers, students, class;
SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE teachers (
    teacher_email VARCHAR(255) PRIMARY KEY
);

CREATE TABLE students (
    student_email VARCHAR(255) PRIMARY KEY,
    is_suspended  BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE class (
    teacher_email VARCHAR(255) NOT NULL REFERENCES teachers (teacher_email),
    student_email VARCHAR(255) NOT NULL REFERENCES students (student_email),
    PRIMARY KEY (teacher_email, student_email)
);

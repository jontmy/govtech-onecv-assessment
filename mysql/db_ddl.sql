SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS teachers, students, teaches;
SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE teachers (
    id    INT          NOT NULL AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL UNIQUE,
    PRIMARY KEY (id)
);

CREATE TABLE students (
    id           INT          NOT NULL AUTO_INCREMENT,
    email        VARCHAR(255) NOT NULL UNIQUE,
    is_suspended BOOLEAN      NOT NULL DEFAULT FALSE,
    PRIMARY KEY (id)
);

CREATE TABLE teaches (
    teacher_id INT NOT NULL,
    student_id INT NOT NULL,
    PRIMARY KEY (teacher_id, student_id),
    FOREIGN KEY (teacher_id) REFERENCES teachers (id),
    FOREIGN KEY (student_id) REFERENCES students (id)
);

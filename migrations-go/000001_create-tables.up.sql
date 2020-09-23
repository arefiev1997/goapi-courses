CREATE TABLE IF NOT EXISTS class (
    id BIGSERIAL PRIMARY KEY,
    number INT NOT NULL,
    letter CHAR(1) NOT NULL
);

CREATE TABLE IF NOT EXISTS student (
    id BIGSERIAL PRIMARY KEY,
    surname TEXT NOT NULL,
    name TEXT NOT NULL,
    patronymic TEXT NOT NULL,
    age INT NOT NULL,
    class BIGINT REFERENCES class (id)
);

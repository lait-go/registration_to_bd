CREATE TABLE IF NOT EXISTS person (
  id serial primary key,
  name varchar NOT NULL,
  surname varchar NOT NULL,
  patronymic varchar,
  age integer NOT NULL,
  gender varchar NOT NULL,
  nationality varchar NOT NULL
);


CREATE TABLE IF NOT EXISTS person (
  id serial primary key,
  person_name varchar NOT NULL,
  password varchar NOT NULL,
  registration_date TIMESTAMP,
  status varchar NOT NULL,
  role varchar DEFAULT 'user' NOT NULL
);


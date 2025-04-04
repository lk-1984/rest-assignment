#!/bin/bash
DB_NAME="atlas"
DB_EXISTS=$(createdb -U postgres ${DB_NAME} 2>/dev/null; echo $?)
echo "Database ${DB_NAME} createdb result: ${DB_EXISTS}"

if [ "${DB_EXISTS}" == 0 ]; then
  psql -U postgres -tAc "CREATE USER api WITH PASSWORD 'api'"
  psql -U postgres -tAc "GRANT ALL PRIVILEGES ON DATABASE atlas TO api;"

  psql -U postgres atlas -tAc "GRANT ALL PRIVILEGES ON SCHEMA public TO api"
  psql -U postgres atlas -tAc "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO api"
  psql -U postgres atlas -tAc "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO api"

  psql -U postgres atlas -tAc "CREATE TABLE continents (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL)"
  psql -U postgres atlas -tAc "CREATE TABLE countries (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, continent_id INT NOT NULL, FOREIGN KEY (continent_id) REFERENCES continents(id))"
  psql -U postgres atlas -tAc "CREATE TABLE cities (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, country_id INT NOT NULL, FOREIGN KEY (country_id) REFERENCES countries(id))"

  psql -U postgres atlas -tAc "GRANT USAGE, SELECT ON SEQUENCE continents_id_seq TO api"
  psql -U postgres atlas -tAc "GRANT USAGE, SELECT ON SEQUENCE countries_id_seq TO api"
  psql -U postgres atlas -tAc "GRANT USAGE, SELECT ON SEQUENCE cities_id_seq TO api"
fi

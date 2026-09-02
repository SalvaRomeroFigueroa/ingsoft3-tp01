-- Este archivo es de referencia / documentación.
-- La app corre estas mismas sentencias automáticamente al arrancar
-- (ver Store.Migrate en store.go), así que no hace falta ejecutarlo a mano.
-- Sirve además como script de inicialización si más adelante lo montan
-- en /docker-entrypoint-initdb.d/ del contenedor de Postgres.

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS lists (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    list_id INTEGER NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pendiente',
    due_date DATE
);

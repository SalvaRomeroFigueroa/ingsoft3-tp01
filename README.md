# To-Do List — backend Go + frontend vanilla JS

App de listas de tareas por usuario.

## Estructura

todo-app/
├── backend/ # API REST en Go (stdlib + driver de Postgres)
│ ├── main.go
│ ├── store.go # acceso a datos
│ ├── rules.go # reglas de negocio puras (testeables sin DB)
│ ├── rules_test.go # tests de ejemplo de esas reglas
│ ├── handlers_*.go # rutas HTTP
│ ├── models.go
│ ├── errors.go
│ ├── http_helpers.go
│ ├── schema.sql # documental, la app migra sola al arrancar
│ └── go.mod
└── frontend/ # SPA sin frameworks (HTML + CSS + JS vanilla)
├── index.html
├── style.css
└── app.js


## Reglas de negocio

Viven en `backend/rules.go`, separadas del acceso a datos para poder
testearlas sin levantar Postgres:

| Regla | Función | Tipo |
|---|---|---|
| Título no vacío, máx. 100 caracteres | `ValidateTaskTitle` | Validación |
| Fecha de vencimiento no puede ser pasada | `ValidateDueDate` | Validación |
| `pendiente ⇄ completada`, ambas → `archivada`; `archivada` es terminal | `ValidateTransition` | Transición de estado |
| No se borra una lista con tareas pendientes | `CanDeleteList` | Restricción |
| Contador de pendientes por lista | `ListsByUser` (en `store.go`) | Cálculo |
| Un usuario solo ve/modifica sus propias listas y tareas | chequeo de `owner` en los handlers | Autorización |

## Cómo correrla en local

Go 1.22+ y PostgreSQL

1. Crear la base:
```bash
   createdb todoapp
```
2. Backend:
```bash
   cd backend
   go mod tidy      # descarga el driver de Postgres (github.com/lib/pq)
   go run .
```
   Por defecto usa `postgres://postgres:postgres@localhost:5432/todoapp?sslmode=disable`.
   Para usar otra, setear `DATABASE_URL` antes de `go run .`. El puerto HTTP
   también es configurable con `PORT` (default `8080`).
3. Frontend: abrir `frontend/index.html` directamente en el navegador (usa
   `http://localhost:8080` como backend, ya está seteado en el HTML). Si
   se prefiere servirlo con un servidor estático:
```bash
   cd frontend
   python3 -m http.server 5500
```

## Cómo levantarla con Docker (backend + frontend + Postgres)

Prerrequisito: Docker Desktop instalado y corriendo (`docker ps` no debe dar error de conexión).

```bash
cp .env.example .env
# editar .env y poner la contraseña que se prefiera para DB_PASSWORD

docker compose up -d --build
docker compose ps      # esperar a ver "db" healthy y backend/frontend "running"
```

Con eso arriba:
- Frontend: http://localhost:3000
- Backend (para probar con curl/Postman): http://localhost:8080/api/health

El frontend le pega a `/api/...` con ruta relativa, y `nginx.conf` la reenvía
internamente al contenedor del backend — por eso, dentro de Docker, no hace
falta tocar `window.API_BASE` en `index.html` (queda vacío por defecto).

Para bajar el sistema conservando los datos:
```bash
docker compose down
```

Para bajar el sistema **y borrar también los datos** (vuelve a un estado limpio, sin usuarios ni listas):
```bash
docker compose down -v
```

Los datos de Postgres viven en el volumen nombrado `db_data`, no en los
contenedores — por eso sobreviven a un `docker compose down`/`up` normal,
y solo se pierden si se agrega el flag `-v`.
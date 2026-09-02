# RESPUESTAS

## TP1

### 1. el conflicto ocurre porque las 2 ramas independientes entre si (las 2 creadas a partir de un main igual), generan un cambio en el mismo archivo, por lo que GIT no puede detectar automáticamente cual conservar, pero si la rama titulo-b hubiera sido creada desde la rama titulo-a entonces no habría problema y se lo interpretaría como una actualización, o tambien si afectaran lineas distintas del archivo

### 2. La unica parte donde tuve problemas fue en el 4.6 para crear el pull request desde GitBash, particularmente para crear el PR, y tambien cuando se creo la rama no se me habían guardado los cambios correctamente y tuve que empezar la parte local del ejercicio de nuevo

### 3. Solo se utilizo IA en la parte del punto 4.6 donde pedía crear un PR desde gitbash, ya que lograba crear la rama titulo-b pero no el PR


## TP2 — Contenedores

### 1. Elección de la app

Elegí construir mi propia to-do list (frontend + backend + base de datos)
en vez de tomar una app de GitHub, para poder armarla desde el principio
cumpliendo los cinco criterios de `elegir-app.md`:

1. **Que pueda ejecutarla hoy**: la escribí yo, la corrí local con
   `go run .` + Postgres antes de este TP.
2. **Comandos de compilación conocidos**: `go build` para el backend; el
   frontend no tiene build (vanilla JS, sin bundler).
3. **Configuración de la base por variable de entorno**: el backend lee
   `DATABASE_URL` con `os.Getenv`, nada hardcodeado.
4. **Lógica para testear**: reglas de negocio separadas en `rules.go`
   (validación de título, validación de fecha, transición de estado,
   restricción de borrado de lista).
5. **Que la entienda para modificarla**: la escribí completa.

### 2. Decisiones de contenerización

**Imágenes base:** `golang:1.22-alpine` para compilar el backend (etapa
`build`) y `alpine:3.20` para correrlo (etapa final) — multi-stage. El
frontend usa `nginx:alpine`, una sola etapa.

**Por qué el frontend NO es multi-stage:** mi frontend es HTML/CSS/JS
vanilla, sin `package.json` ni bundler — no hay ningún paso de compilación
que separar del runtime. Los archivos que sirve nginx son exactamente los
del repo. 

**Ruta del frontend al backend:** el frontend llama a `/api/...` con ruta
relativa, y `nginx.conf` la reenvía a `http://backend:8080` (nombre del
servicio en compose). La elegí en vez de que el navegador llame directo al
backend porque evita configurar CORS y la misma imagen del frontend sirve
en cualquier entorno sin URLs hardcodeadas.

**Qué persiste:** solo los datos de Postgres, en el volumen nombrado
`db_data`. Backend y frontend son stateless.

**Healthcheck de la base:** `pg_isready -U postgres`, con
`condition: service_healthy` en el `depends_on` del backend.
`depends_on` a secas solo garantiza el orden de arranque, no que Postgres
ya acepte conexiones — sin el healthcheck, el backend podría arrancar y
fallar su primer `Ping()` porque la base todavía está inicializando.

**Secretos:** la contraseña de Postgres viaja por `${DB_PASSWORD}`, leída
de un `.env` que no se commitea (sí se commitea `.env.example`, con un
valor de ejemplo, sin la contraseña real).

**Arquitectura:** construida en windows/amd64.

### 3. Problemas encontrados y cómo los resolví

**Problema 1 — Puerto 8080 ya ocupado.**
Al levantar todo con `docker compose up -d --build`, el backend no arrancó:
`Bind for 0.0.0.0:8080 failed: port is already allocated`. La causa: un
contenedor de prueba (`backend-test`) que había levantado antes, para
probar el backend suelto contra mi Postgres local, había quedado corriendo
y ocupando ese puerto. Lo resolví con `docker rm -f backend-test` antes de
volver a levantar el compose completo.

### 4. Verificación

Levanté el sistema completo con `docker compose up -d --build` y confirmé
con `docker compose ps` que los tres servicios (`db`, `backend`,
`frontend`) quedaron arriba, con `db` en estado `healthy`. Probé el
endpoint `/api/health` del backend por curl, y usé el frontend en
`http://localhost:3000` para crear un usuario, una lista y una tarea.
Comprobé la persistencia de los datos: tras un `docker compose down` +
`up` (sin `-v`) la lista seguía existiendo; tras `docker compose down -v`
+ `up`, el volumen se recreó vacío y la lista ya no aparecía (aunque al
volver a loguearme con el mismo nombre de usuario, el backend lo creó de
nuevo con id 1, por tratarse de una tabla nueva). También comparé el
tamaño de la imagen final del backend (27.1MB, basada en `alpine:3.20`)
contra el peso de la imagen de build (`golang:1.22-alpine`, mucho más
pesada), confirmando en números el beneficio del multi-stage.

### 5. Uso de IA

Usé Claude para guiarme paso a paso en la escritura de los
Dockerfiles, `docker-compose.yml` y `nginx.conf` de este TP. Escribí cada
archivo y comando de cmd yo mismo siguiendo esa guía, y resolví los tres problemas de la
sección anterior con ayuda de Claude para interpretar los mensajes de error de Docker.

## TP3 - DevOps

### 1. Duración de sprint

Puse que dure 1 semana cada sprint, porque es cada 1 semana que empezamos un nuevo tp normalmente

### 2. - Limite de In Progress

Puse un límite de 2 tareas, debido a la justificación dada en el video, hacer un limite muy alto puede llevar a empezar muchas tareas y no terminarlas, un limite bajo fuerza a terminar tareas antes de asumir otras nuevas

### 3. - Historia mal escrita

Una historia de usuario debería describir una necesidad final de un usuario, no las herramientas básicas para lograrlo, describe funciones en lugar de elementos del sistema (por ejemplo pueden plantear un sistema de reserva de turnos de un consultorio, pero no pueden simplemente pedir una entidad cliente x turno)

Son las tasks las que plantean elementos internos del sistema

### 4. - Problemas

Como seguí el video del tp, mientras leía el documento del mismo, no tuve ningun problema mayor fuera de confusiones en el orden de desarrollo de los pasos por discrepancias entre ambas fuentes

### 5. - Uso de IA

Para este tp no utilicé la ayuda de la IA en ningún paso


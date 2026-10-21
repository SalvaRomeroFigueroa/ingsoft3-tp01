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

## TP4 — CI: Pipelines as Code

### 1. Estructura del pipeline

Dos jobs, `build-backend` y `build-frontend`, uno por cada Dockerfile del
TP2. Corren en paralelo porque son independientes entre sí — ninguno
necesita el resultado del otro, así que serializarlos solo sumaría tiempo
sin ningún beneficio.

### 2. Qué cachea mi pipeline

Las capas de cada imagen Docker, vía `type=gha` con `mode=max`. Cada job
usa un `scope` distinto (`backend` / `frontend`) para no pisarse el cache
entre sí. Si el cache desaparece, el pipeline sigue funcionando igual, solo
más lento: es una optimización, no una dependencia. Si fallara sin él,
sería un bug del Dockerfile, no del cache.

### 3. Por qué construye con mi Dockerfile

El workflow no tiene ninguna línea de Go — solo invoca `docker build` sobre
los Dockerfiles. Así evito tener dos definiciones de build (la del
Dockerfile y la del pipeline) que con el tiempo terminan divergiendo: lo
que se verifica acá es exactamente lo que después se despliega.

### 4. Problemas

Seguí el video explicativo, no tuve problemas mayores mas allá de una pequeña confusión en la parte de hacer 2 PR para probar freno por out of date

### 5. - Uso de IA

Para este tp no utilicé la ayuda de la IA en ningún paso

## Corrección TP2

El profesor marcó dos faltantes en la entrega original: (1) el Dockerfile del
frontend no era multi-stage, (2) faltaba `docker-compose.registry.yml` y las
imágenes nunca se publicaron de verdad en el registry.

Corregido así: el frontend pasó a tener dos etapas — la primera copia los
estáticos y corre `nginx -t` para validar la sintaxis de `nginx.conf` antes
de hornear la imagen final (comprobé que la etapa frena de verdad: rompí el
`nginx.conf` a propósito y el build falló ahí). Las imágenes se publicaron en
ghcr.io como `ghcr.io/salvaromerofigueroa/todo-backend:v0.1.0` y
`todo-frontend:v0.1.0`, públicas, y probadas bajándolas sin estar logueado
con `docker-compose.registry.yml`.

## TP5 — Testing y calidad de código

### 1. Qué lógica elegí testear y por qué esa

Elegí testear donde un bug le hace daño real a un usuario:

- **Autorización por dueño** (handlers del backend): si falla, un usuario ve, modifica o borra las listas de otro. Es el bug más caro de la app.
- **Reglas de negocio de `rules.go`**: validación del título y de la fecha, transiciones de estado (`archivada` es terminal) y la restricción de no borrar una lista con tareas pendientes. Un bug acá deja datos inválidos o pierde tareas.
- **Frontend (`lib/tareas.js` y `lib/api.js`)**: qué acciones ofrece la UI según el estado de la tarea (que no ofrezca una transición que el backend rechaza) y el manejo de errores del cliente de la API (que el usuario vea por qué falló un pedido).

No testeé el SQL de `store.go` ni el DOM de `app.js`: necesitan tests de otro tipo (integración con Postgres real, o de interfaz), que quedan fuera de este TP.

### 2. Herramientas usadas por fila (mi stack es Go + vitest, no .NET)

| Pieza | Backend (Go) | Frontend (vitest) |
|---|---|---|
| Test parametrizado | Tests por tabla con `t.Run` (`TestDeleteList_Autorizacion`, `TestCreateTask_Validaciones`, `TestUpdateTaskStatus_Transiciones`) | `it.each` (`validarTitulo`, `accionesPara`) |
| Caso de error | `TestPostLists_NombreVacio_Rechazado`, `TestPostLists_ErrorDelStore_Devuelve500` | Rechazo con el mensaje del servidor y con mensaje genérico en `api.test.js` |
| Doble / mock | `FakeStore` hecho a mano (sin librería), con contador de llamadas | `vi.fn()` en lugar de `fetch` |
| Medidor de cobertura | `go test -coverprofile` + `go tool cover` | vitest con proveedor v8 |
| Umbral que frena | Script de bash en el pipeline (Go no trae una bandera para esto) | `thresholds` en `vitest.config.js` |
| Qué entra en la cuenta | `grep -v` sobre el perfil de cobertura | `coverage.include: ["lib/**"]` |

Los tests corren dentro de la etapa `test` de cada Dockerfile (`FROM build AS test`), no con herramientas nativas en el runner. Así se prueba el mismo código y las mismas dependencias que se construyen para producción, sin una segunda receta de build.

### 3. Refactor para poder mockear

**Backend.** Antes, los handlers recibían `*Store`, el struct concreto atado a `*sql.DB`: no había forma de reemplazar la base en un test, y cada test hubiera necesitado Postgres levantado. Extraje la interfaz `DataStore` (`datastore.go`) con los métodos que usan los handlers y cambié la firma de los tres `register...Routes` de `*Store` a `DataStore`. `Store` la cumple sin cambios (las interfaces de Go se satisfacen estructuralmente) y `main.go` no se tocó. Los tests inyectan `FakeStore`, que guarda todo en mapas en memoria y cuenta las llamadas a `CreateList`: `TestPostLists_CreaListaConDoble` verifica la respuesta HTTP y que `CreateList` se llamó exactamente una vez (verificación de interacción, no solo de resultado).

**Frontend.** `app.js` llamaba al `fetch` global y leía el `currentUser` global, todo mezclado con el DOM. Extraje `crearClienteApi({ fetchFn, baseUrl, getUserId })` a `lib/api.js`: recibe `fetch` desde afuera, así el test le pasa un `vi.fn()` y verifica a qué URL llamó, cuántas veces y con qué header `X-User-Id`. La lógica pura salió a `lib/tareas.js`.

Dos consecuencias: `index.html` carga `app.js` como módulo ES, así que ya no se puede abrir con doble clic (`file://`); sin Docker hay que servir la carpeta (`python -m http.server`). Y los tests están en `frontend/test/`, no junto al código, para que el `COPY lib` del Dockerfile no lleve tests a la imagen de producción.

### 4. Umbral de cobertura

**Backend: 75 %, sobre sentencias** (es la única métrica que mide Go). Cuando lo fijé medí 82,5 % en mi máquina y 81,9 % en la etapa de tests del contenedor, ambos sin `main.go` ni `store.go`. Elegí 75 porque queda unos 7 puntos por debajo de lo que medía: arranca en verde y no se rompe por ruido (el contenedor mide algo menos que mi máquina), pero código nuevo sin tests lo hace caer. Si lo subiera diez puntos (85) quedaría por encima de lo que mido, y el pipeline fallaría desde el primer commit sin que haya ninguna regresión real. Con la suite completa, el pipeline mide hoy **84,3 %**.

**Frontend: 90 %, en líneas y en ramas.** Medí 100 % en todo `lib/`. Puse las dos métricas porque según la versión de vitest una función nueva sin tests puede no caer en ramas; en mi versión (5.0.5) cayó en las dos. Margen de 10 puntos, y como `lib/` es chica, una función nueva sin tests ya hace caer el número.

**Cobertura de rama hoy:** frontend 100 %. En el backend no tengo el dato: Go mide sentencias, no ramas.

El resumen de cobertura de los dos lados y los reportes descargables (artefactos `coverage-backend` y `coverage-frontend`) están en esta corrida: https://github.com/SalvaRomeroFigueroa/ingsoft3-tp01/actions/runs/39545995035

### 5. Qué dejé afuera de la cuenta de cobertura

- **Backend:** `main.go` (el arranque: lee variables de entorno, abre la conexión y registra rutas, sin reglas de negocio) y `store.go` (queries SQL contra Postgres: los tests unitarios lo reemplazan con el `FakeStore`, y cubrirlo pide tests de integración con una base real, que no hice). `models.go` no aparece en el perfil: son structs sin sentencias ejecutables. Es el mismo filtro para el umbral y para el reporte, para no mostrar un número y exigir otro. Lo escribí como lista de lo que se saca, no de lo que se deja, así un `.go` nuevo entra a la cuenta solo.
- **Frontend:** entra solo `lib/**` (`coverage.include`). `app.js` es pegamento de DOM y queda afuera. Limitación: acá la lista es de lo que se incluye, así que código nuevo fuera de `lib/` no se mide. Lo compenso poniendo la lógica nueva en `lib/` (el PR abierto lo mostró: un archivo nuevo en `lib/` entró solo a la cuenta).

### 6. Por qué coverage alto no garantiza calidad

Un test que llame a `accionesPara("pendiente")` sin ningún `expect` ejecuta la función y suma el 100 % de sus líneas, pero no verifica nada: si alguien cambiara el `case "pendiente"` para ofrecer "reabrir", ese test seguiría en verde. En el backend, `ValidateTransition` figuraba como cubierta en su línea `if !ok || !allowed[next]` aunque el lado `!ok` nunca se evaluaba como verdadero: las sentencias miden que la línea corrió, no que cada condición haya tomado cada valor. Y mi 100 % del frontend es solo de `lib/`: no dice nada sobre si la UI (`app.js`) funciona.

### 7. El ejercicio del camino sin cubrir

- **Qué línea es:** en `rules.go`, `ValidateTransition`, el `return nil` dentro de `if current == next` (pedir el estado que la tarea ya tiene). Ningún test lo recorría; el reporte marcaba la función en 83,3 %.
- **Qué entrada la recorre:** `ValidateTransition("pendiente", "pendiente")`.
- **Qué decidí:** lo agregué (`TestValidateTransition_MismoEstado_EsNoOp`), porque define un comportamiento real: pedir el estado actual es un no-op aceptado, no un error. El total filtrado subió de 82,1 % a 82,5 %. Sumé también `TestValidateTransition_EstadoActualDesconocido` para el lado `!ok` de la condición, que las sentencias daban por cubierto sin estarlo.

### 8. Pull Requests bloqueados por cobertura

**PR mergeado (backend): la secuencia completa.** Agregué en `rules.go` `PrioridadDeTarea`, `ResumenDeLista` y `Etiqueta`: unas 35 sentencias nuevas sin tests. Compilaba y los tests existentes pasaban, pero el check `build-backend` se puso en **rojo** en el paso "Resumen y umbral de coverage del backend": `Cobertura de statements: 70.7% (umbral 75%)`. El paso de tests estaba en verde: frenó por calidad, no por compilación. Métrica: sentencias.

Para arreglarlo escribí `rules_prioridad_test.go`, un caso por cada salida que el código declara: 10 casos de tabla para `PrioridadDeTarea` (con los bordes de 7/8 y 30/31 días, y que la hora del día no cambia el resultado), 3 tests de `ResumenDeLista` (cuenta por estado, sin tareas activas, todo completado) y 4 casos de `Etiqueta`. El pipeline pasó de 70,7 % a 84,3 %, el check quedó en verde y se mergeó. Ese código quedó en `main` como código de la demostración; ningún endpoint lo usa.

- Corrida roja por umbral: https://github.com/SalvaRomeroFigueroa/ingsoft3-tp01/actions/runs/39517802896
- Corrida verde tras los tests: https://github.com/SalvaRomeroFigueroa/ingsoft3-tp01/actions/runs/39545995035
- Pull Request mergeado (la conversación muestra la secuencia): https://github.com/SalvaRomeroFigueroa/ingsoft3-tp01/pull/21

**PR abierto (frontend): el freno vigente.** Agregué `frontend/lib/prioridad.js` sin tests y dejé el PR abierto y sin arreglar. `build-frontend` está en **rojo** por cobertura, en las dos métricas:

- `ERROR: Coverage for lines (54.9%) does not meet global threshold (90%)`
- `ERROR: Coverage for branches (59.09%) does not meet global threshold (90%)`

Frenó por líneas y por ramas, con vitest 5.0.5. Los tests pasan; falla solo el umbral.

- Pull Request abierto: https://github.com/SalvaRomeroFigueroa/ingsoft3-tp01/pull/22
- Corrida roja: https://github.com/SalvaRomeroFigueroa/ingsoft3-tp01/actions/runs/39549651844

### 9. Problemas encontrados y cómo los resolví

- **PowerShell vs. Git Bash.** `tail` y `grep` no existen en PowerShell, y `go tool cover -func=coverage.out` fallaba con `too many arguments` porque PowerShell parte el argumento en el punto (y dejó un archivo `coverage` suelto). Lo resolví entrecomillando el argumento y usando Git Bash para los comandos con pipes.
- **`go test` mostraba solo 7 tests.** Solo existían 3 de los 8 archivos de test; los otros no se habían creado. Lo detecté con `ls *_test.go` y repetí la medición con `-count=1` para evitar el caché.
- **Un 17 % de cobertura al arrancar.** El total incluía `store.go` y `main.go`, ambos en 0 %. Lo resolví con el filtro de lo que no se testea de forma unitaria (sección 5).
- **`git add -A` iba a subir basura.** Entraban reportes generados (`frontend/coverage`, `frontend-coverage`) y un `package.json` que había creado por error dentro de `backend/`. Lo vi en el `git status` antes de commitear: `git reset`, borré lo que sobraba y agregué los patrones al `.gitignore`.
- **82,5 % en mi máquina, 81,9 % en el contenedor.** Decidí que el número de referencia es el del contenedor y el pipeline, porque es el que se mide en cada PR.

### 10. Uso de IA

Usé Claude para diseñar el refactor a `DataStore`, escribir el `FakeStore`, los tests del backend y del frontend, la etapa `test` de los Dockerfiles y los pasos nuevos de `ci.yml`. Lo verifiqué corriendo ambas etapas de test en mi máquina con `docker build --target test` y `docker run`, midiendo la cobertura localmente y viendo las corridas en Actions, incluida la corrida que falló a propósito. Puedo explicar qué verifica cada assert. Lo que **no** está cubierto: el SQL de `store.go`, `main.go` y la UI de `app.js`, que quedan afuera de la cuenta a propósito; y, dentro de la cuenta: Lo que **no** está cubierto, según el reporte de cobertura: afuera de la cuenta, el SQL de `store.go`, `main.go` y la UI de `app.js`; y dentro de la cuenta, 20 bloques de `handlers_lists.go`, `handlers_tasks.go` y `handlers_users.go` (39 sentencias). De esos, 19 son ramas de error defensivas y uno es un camino feliz:

- **El camino feliz de `GET /api/lists`** (el `writeJSON` con 200, `handlers_lists.go:23`): ningún test lo recorre. Es el hueco que más me preocupa, porque es la ruta que más usa el frontend.
- **401 por falta de usuario** en `POST /api/lists`, `DELETE /api/lists/{id}`, el `GET` de tareas y el `DELETE` de tareas. El 401 sí está cubierto en `GET /api/lists`, en el alta de tareas y en el cambio de estado.
- **400 por id inválido** en las cuatro rutas de tareas, y **400 por cuerpo JSON inválido** en `POST /api/lists`, en el alta de tareas y en el cambio de estado.
- **404 por lista inexistente** en el `GET` de tareas.
- **500 cuando falla la base**, en siete lugares: `PendingTaskCount` y `DeleteList` (borrado de listas), `TasksByList`, `CreateTask`, `UpdateTaskStatus`, `DeleteTask` y `FindOrCreateUser` (login). Solo simulé la caída de la base en `GET` y `POST` de listas, con un store que falla a propósito.

Dejé sin cubrir los 500 porque cada uno pide un doble que falle en un método distinto, y su valor es bajo frente al riesgo de las reglas de negocio y la autorización, que sí están cubiertas. El camino feliz de `GET /api/lists` es la excepción: ahí el hueco es real, y lo declaro en vez de esconderlo.
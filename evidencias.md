# Evidencias — TP1

## 1. Push directo a main rechazado

<img width="890" height="453" alt="4 4" src="https://github.com/user-attachments/assets/9e17f351-6ed3-4b20-944e-f425ab9eb555" />

GitHub rechaza el push porque main está protegida y la regla alcanza también al dueño del repo.

## 2. El PR de la rama B no se puede mergear: conflicto

<img width="1920" height="974" alt="4 6 1" src="https://github.com/user-attachments/assets/52664ae9-fb62-4e48-93ce-b7c9a79dacfa" />

Conflictos entre la version A (ya mergeada), con la version B (en pull request)

## 3. Linea en conflicto en el README.md

<img width="1920" height="968" alt="4 6 2" src="https://github.com/user-attachments/assets/6d4874b1-b3b1-4e29-9fd0-2504b37ebb87" />

Se muestra la linea en conflicto del README, y que corresponde a cada rama

## 4. Release v1.0.0 publicada

<img width="1917" height="967" alt="4 7" src="https://github.com/user-attachments/assets/654a0d44-6134-40e6-85b4-813f8b55c986" />
Se muestra la release v1.0.0 publicada en GitHub, asociada al tag correspondiente y con las notas de la primera versión estable del TP.

# Evidencias - TP2 — Contenedores

## 1. Build de la imagen del backend (multi-stage)

<img width="1920" height="1020" alt="1 - Dockerizacion backend - 1" src="https://github.com/user-attachments/assets/7f248c5b-7154-4dfa-a60c-92cee54a912e" />
<img width="1920" height="1020" alt="1 - Dockerizacion backend - 2" src="https://github.com/user-attachments/assets/a05eff0a-ff63-4ecf-8f7e-9348f6a593e3" />

Se muestra el build exitoso del Dockerfile del backend, con las dos etapas (`build` con el compilador de Go y `stage-1` con la imagen final de Alpine).

## 2. Contenedor del backend respondiendo /health

<img width="1265" height="163" alt="2 - testear backend 1" src="https://github.com/user-attachments/assets/6e3a74f1-6a7f-460a-92c9-0f6c8609350d" />
<img width="1407" height="85" alt="2 - testear backend 2" src="https://github.com/user-attachments/assets/d11f8b40-16c3-4e8e-b66b-c8cb67d64e7a" />

El contenedor del backend, corrido de forma aislada y conectado a Postgres, responde correctamente al endpoint de salud.

## 3. Los tres servicios levantados con docker compose

<img width="1920" height="1020" alt="3 - Docker build - 1" src="https://github.com/user-attachments/assets/c355f408-f9f7-403b-8458-ba95b01be950" />
<img width="1920" height="1020" alt="3 - Docker build - 2" src="https://github.com/user-attachments/assets/95b65d6f-b46f-42c6-80fe-70d19733cc5d" />
<img width="1920" height="1020" alt="3 - Docker build - 3" src="https://github.com/user-attachments/assets/5e8bc460-8931-452d-b077-42e46cd1752f" />
<img width="1917" height="160" alt="3 - Docker build - 4 - compose" src="https://github.com/user-attachments/assets/69298df4-0286-4a32-b39e-9faf8bfe3fcd" />

Salida de `docker compose ps` mostrando la base de datos en estado healthy y los contenedores de backend y frontend corriendo.

## 4. Frontend funcionando end-to-end

<img width="1920" height="1020" alt="4 - Página dockerizada andando" src="https://github.com/user-attachments/assets/20a1b838-37b4-4e79-943e-0fb20c154983" />

La aplicación corriendo en el navegador, con una lista y una tarea creadas a través del frontend, que llegan al backend por el proxy de nginx.

## 5. Persistencia de datos tras `docker compose down` / `up`

<img width="1911" height="591" alt="5 - persistencia" src="https://github.com/user-attachments/assets/3c191865-5a26-4446-815e-3beb75ff7cc0" />

Después de bajar y volver a levantar los contenedores sin borrar el volumen, la lista creada previamente sigue existiendo: los datos persisten porque viven en el volumen `db_data`, no en los contenedores.

## 6. Pérdida de datos tras `docker compose down -v` / `up`

<img width="1916" height="587" alt="6 - no persistencia de volumen" src="https://github.com/user-attachments/assets/59c5a0dc-51e9-4f00-87e2-1c1d7764990f" />

Al bajar los contenedores borrando también el volumen (`-v`) y volver a levantar, la base de datos arranca vacía: la lista ya no aparece, lo que confirma que el volumen era efectivamente lo único que sostenía los datos.

## 7. Comparación de tamaño de imágenes

<img width="1123" height="672" alt="7 - pesos imagenes compilador go vs final" src="https://github.com/user-attachments/assets/d01e7178-49dd-45b9-8806-02e5853c4af2" />

Comparación entre el tamaño de la imagen de build (`golang:1.22-alpine`) y la imagen final del backend (basada en `alpine:3.20`), que evidencia el ahorro de espacio logrado con el enfoque multi-stage.
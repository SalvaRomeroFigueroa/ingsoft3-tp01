# ---- Etapa 1: validación ----
# No hay bundler (HTML/CSS/JS vanilla, sin build step) — pero esta etapa sí
# cumple una función real: valida la sintaxis de nginx.conf ANTES de hornear
# la imagen final. Si el archivo tiene un error, el build FALLA ACÁ, no en
# producción cuando el contenedor intente arrancar.
FROM nginx:alpine AS build
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY index.html style.css app.js /usr/share/nginx/html/
RUN nginx -t

# ---- Etapa 2: runtime (imagen final, ya validada) ----
FROM nginx:alpine
COPY --from=build /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf
COPY --from=build /usr/share/nginx/html /usr/share/nginx/html
EXPOSE 80
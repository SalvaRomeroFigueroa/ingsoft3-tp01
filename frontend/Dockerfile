# ---- Etapa 1: build (Node + dependencias de desarrollo) ----
# No hay bundler: esta etapa instala las dependencias de test y junta el
# código fuente. Lo que sale de acá hacia producción es solo lo que la
# etapa final copia explícitamente (sin node_modules ni tests).
FROM node:22-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .

# ---- Etapa 2: tests (reutiliza build: mismo código y mismas dependencias) ----
FROM build AS test
ENTRYPOINT ["npm", "run", "test:ci"]

# ---- Etapa 3: runtime (nginx sirve los estáticos) ----
FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/index.html /app/style.css /app/app.js /usr/share/nginx/html/
COPY --from=build /app/lib /usr/share/nginx/html/lib
RUN nginx -t
EXPOSE 80
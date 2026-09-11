# ==============================================================================
# Dockerfile Unificado para Despliegue en la Nube (Hugging Face Spaces - Port 7860)
# ==============================================================================

# Etapa 1: Compilar el microservicio en Go
FROM golang:1.21-alpine AS go-builder
WORKDIR /app/api-go
RUN apk add --no-cache git
COPY api-go/go.mod api-go/go.sum ./
RUN go mod download
COPY api-go/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app/api-go/server main.go

# Etapa 2: Preparar dependencias de Node.js
FROM node:18-alpine AS node-builder
WORKDIR /app/api-node
COPY api-node/package*.json ./
RUN npm ci --only=production || npm install --omit=dev
COPY api-node/ .

# Etapa 3: Imagen unificada de producción
FROM node:18-alpine

RUN apk add --no-cache nginx bash ca-certificates tzdata && \
    mkdir -p /run/nginx /usr/share/nginx/html /app/api-go /app/api-node

# Copiar microservicio Go compilado
COPY --from=go-builder /app/api-go/server /app/api-go/server

# Copiar microservicio Node.js
COPY --from=node-builder /app/api-node /app/api-node

# Copiar Frontend estático y configuración de Nginx
COPY frontend/index.html /usr/share/nginx/html/index.html
COPY nginx-cloud.conf /etc/nginx/http.d/default.conf

# Copiar y dar permisos al script de entrada
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Puerto web expuesto
EXPOSE 8080

CMD ["/app/entrypoint.sh"]

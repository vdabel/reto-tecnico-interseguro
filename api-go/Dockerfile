# Etapa 1: Construcción
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar git y certificados necesarios para descargar dependencias
RUN apk add --no-cache git

# Copiar manifiestos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente y compilar binario optimizado
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o server main.go

# Etapa 2: Imagen mínima de producción
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 3000

ENV PORT=3000

CMD ["./server"]

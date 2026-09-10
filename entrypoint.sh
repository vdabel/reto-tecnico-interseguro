#!/bin/bash
set -e

echo "=== Iniciando Ecosistema de Microservicios en la Nube ==="

# 1. Iniciar Microservicio Node.js en segundo plano (puerto 4000)
echo "-> Arrancando API Node.js (Estadisticas) en puerto 4000..."
PORT=4000 JWT_SECRET=interseguro_super_secret_key_2024 NODE_ENV=production node /app/api-node/index.js &

# 2. Iniciar Microservicio Go en segundo plano (puerto 3000)
echo "-> Arrancando API Go (Fiber + Gonum QR) en puerto 3000..."
PORT=3000 NODE_API_URL=http://127.0.0.1:4000 JWT_SECRET=interseguro_super_secret_key_2024 /app/api-go/server &

# Esperar 2 segundos para inicializacion interna
sleep 2

# 3. Iniciar Nginx en primer plano (puerto 7860 para Hugging Face)
echo "-> Arrancando Nginx (Frontend + Reverse Proxy) en puerto 7860..."
nginx -g "daemon off;"

FROM node:18-alpine

WORKDIR /app

# Copiar archivos de dependencias
COPY package*.json ./

# Instalar dependencias de producción
RUN npm ci --only=production || npm install --omit=dev

# Copiar el resto del código
COPY . .

EXPOSE 4000

ENV PORT=4000
ENV NODE_ENV=production

CMD ["node", "index.js"]

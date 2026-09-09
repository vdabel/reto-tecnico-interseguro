FROM nginx:alpine

# Reemplazar configuración default de nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copiar aplicación web
COPY index.html /usr/share/nginx/html/index.html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

FROM node:22-alpine AS builder

WORKDIR /app

COPY src/package*.json ./
RUN npm ci
COPY src/ .
RUN npm run build -- --configuration production


FROM nginx:alpine

# Remove the default Nginx static assets
RUN rm -rf /usr/share/nginx/html/*

COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist/frontend/browser /usr/share/nginx/html
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
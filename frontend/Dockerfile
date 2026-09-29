FROM node:24-bookworm-slim AS builder

WORKDIR /src

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:1.28-alpine

COPY --from=builder /src/dist /usr/share/nginx/html
COPY nginx/templates/default.conf.template /etc/nginx/templates/default.conf.template

EXPOSE 80

FROM node:alpine AS builder

ARG VITE_API_URL

WORKDIR /app

COPY package.json package.json
COPY yarn.lock yarn.lock
RUN yarn install

COPY ./ ./
RUN VITE_API_URL=${VITE_API_URL} yarn build


FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /app/dist /usr/share/nginx/html

# === Build Stage ===
FROM node:20 AS build
WORKDIR /app
ARG BUILD_ENV=production
COPY package*.json ./
RUN npm install
# 環境ごとのenvファイルをコピー
COPY . .
RUN npm run build -- --mode ${BUILD_ENV}

# === Production Stage ===
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

# 构建 Vue/Vite 静态资源，Nginx 托管 + 反代 /api → Gin
FROM node:20-alpine AS build
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine
# 与后端在同一 Docker 网络时，compose 服务名一般为 app；后端在宿主机则用 http://host.docker.internal:3000
ARG API_UPSTREAM=http://app:3000
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx/default.conf /etc/nginx/conf.d/default.conf
RUN sed -i "s|__API_UPSTREAM__|${API_UPSTREAM}|g" /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

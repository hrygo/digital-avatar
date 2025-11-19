# 多阶段构建 - TwinOS 智能决策副驾系统
# Build Stage - 构建应用
FROM golang:1.21-alpine AS backend-builder

# 设置工作目录
WORKDIR /app/backend

# 安装构建依赖
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# 复制 go mod 文件
COPY backend/go.mod backend/go.sum ./

# 下载依赖
RUN go mod download && go mod tidy

# 复制源代码
COPY backend/ .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o twin-os-server .

# Frontend Build Stage - 构建前端应用
FROM node:18-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app/frontend

# 复制 package 文件
COPY frontend/package*.json ./

# 安装依赖
RUN npm ci --only=production

# 复制源代码
COPY frontend/ .

# 构建前端应用
RUN npm run build

# Final Stage - 生产镜像
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite tzdata curl

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 -S twinos && \
    adduser -u 1001 -S twinos -G twinos

# 创建应用目录
RUN mkdir -p /app/data /app/backups /app/logs /app/static && \
    chown -R twinos:twinos /app

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件
COPY --from=backend-builder /app/backend/twin-os-server .
COPY --from=frontend-builder /app/frontend/build ./static

# 复制配置文件
COPY backend/.env.example .env

# 设置权限
RUN chmod +x twin-os-server

# 切换到非root用户
USER twinos

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 启动命令
CMD ["./twin-os-server"]
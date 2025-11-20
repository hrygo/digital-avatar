# TwinOS 智能决策副驾系统
# Makefile for building and running TwinOS components

.PHONY: help build clean test dev server web docker install format lint

# 默认目标
.DEFAULT_GOAL := help

# 项目配置
PROJECT_NAME := twin-os
SERVER_NAME := twin-os-server
WEB_NAME := twin-os-web
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# 目录配置
BACKEND_DIR := backend
FRONTEND_DIR := frontend
BUILD_DIR := build
DIST_DIR := dist

# Go配置
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
CGO_ENABLED := 1

# Node.js配置
NODE_VERSION := $(shell node --version 2>/dev/null || echo "unknown")
NPM_VERSION := $(shell npm --version 2>/dev/null || echo "unknown")

help: ## 显示此帮助信息
	@echo 'TwinOS 智能决策副驾系统'
	@echo '=========================='
	@echo ''
	@echo '使用方法:'
	@echo '  make [target]'
	@echo ''
	@echo '可用目标:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean: ## 清理构建文件和缓存
	@echo "🧹 清理构建文件..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@rm -f $(BACKEND_DIR)/$(SERVER_NAME)
	@rm -rf $(FRONTEND_DIR)/build $(FRONTEND_DIR)/dist
	@rm -f $(BACKEND_DIR)/*.log
	@cd $(BACKEND_DIR) && go clean -cache -modcache -testcache
	@cd $(FRONTEND_DIR) && npm run clean 2>/dev/null || true
	@echo "✅ 清理完成"

deps: ## 安装项目依赖
	@echo "📦 安装Go依赖..."
	@cd $(BACKEND_DIR) && go mod download && go mod tidy
	@echo "📦 安装Node.js依赖..."
	@cd $(FRONTEND_DIR) && npm install
	@echo "✅ 依赖安装完成"

setup-db: ## 设置数据库和创建表
	@echo "🗄️ 设置数据库..."
	@mkdir -p data
	@if [ ! -f data/twin-os.db ]; then \
		echo "📝 创建主数据库..."; \
		touch data/twin-os.db; \
	fi
	@echo "🔧 创建微信相关数据表..."
	@if [ -f $(BACKEND_DIR)/create_wechat_tables.sql ]; then \
		sqlite3 data/twin-os.db < $(BACKEND_DIR)/create_wechat_tables.sql; \
		echo "✅ 微信数据表创建完成"; \
	else \
		echo "⚠️ 微信表创建脚本不存在"; \
	fi
	@echo "✅ 数据库设置完成"

setup-wechat: ## 配置微信数据库
	@echo "🔧 配置微信数据库..."
	@if [ -f $(BACKEND_DIR)/test-wechat.db ]; then \
		echo "📱 发现测试微信数据库，正在配置..."; \
		env -u http_proxy -u https_proxy -u all_proxy curl -X PUT http://localhost:1234/api/v1/settings -H "Content-Type: application/json" -d '{"wechat_db_path": "./backend/test-wechat.db"}' 2>/dev/null || echo "⚠️ 请确保后端服务正在运行"; \
		echo "✅ 微信数据库配置完成"; \
	else \
		echo "📝 未找到微信数据库，请手动配置微信数据库路径"; \
		echo "   使用方法: make setup-wechat-db-path DB_PATH=/path/to/wechat.db"; \
	fi

setup-wechat-db-path: ## 设置自定义微信数据库路径 (用法: make setup-wechat-db-path DB_PATH=/path/to/wechat.db)
	@if [ -z "$(DB_PATH)" ]; then \
		echo "❌ 请指定微信数据库路径: make setup-wechat-db-path DB_PATH=/path/to/wechat.db"; \
		exit 1; \
	fi
	@echo "🔧 设置微信数据库路径: $(DB_PATH)"
	@if [ -f "$(DB_PATH)" ]; then \
		env -u http_proxy -u https_proxy -u all_proxy curl -X PUT http://localhost:1234/api/v1/settings -H "Content-Type: application/json" -d '{"wechat_db_path": "$(DB_PATH)"}' || echo "⚠️ 请确保后端服务正在运行"; \
		echo "✅ 微信数据库路径设置完成"; \
	else \
		echo "❌ 微信数据库文件不存在: $(DB_PATH)"; \
	fi

build-server: ## 构建后端服务器
	@echo "🔨 构建TwinOS服务器..."
	@mkdir -p $(BUILD_DIR)
	@cd $(BACKEND_DIR) && go build $(LDFLAGS) -o ../$(BUILD_DIR)/$(SERVER_NAME) .
	@echo "✅ 服务器构建完成: $(BUILD_DIR)/$(SERVER_NAME)"

build-web: ## 构建前端Web应用
	@echo "🔨 构建TwinOS Web应用..."
	@mkdir -p $(DIST_DIR)
	@cd $(FRONTEND_DIR) && npm run build
	@cp -r $(FRONTEND_DIR)/build $(DIST_DIR)/$(WEB_NAME)
	@echo "✅ Web应用构建完成: $(DIST_DIR)/$(WEB_NAME)"

build: build-server build-web ## 构建所有组件

dev-server: ## 运行开发模式服务器
	@echo "🚀 启动TwinOS开发服务器..."
	@cd $(BACKEND_DIR) && go run main.go

dev-web: ## 运行开发模式Web应用
	@echo "🚀 启动TwinOS开发Web应用..."
	@cd $(FRONTEND_DIR) && npm start

dev: ## 同时运行服务器和Web应用(开发模式)
	@echo "🚀 启动TwinOS开发环境..."
	@make -j2 dev-server dev-web

restart: ## 重启开发环境
	@echo "🔄 重启TwinOS开发环境..."
	@echo "🛑 停止现有服务..."
	@pkill -f "go run main.go" 2>/dev/null || true
	@pkill -f "react-scripts start" 2>/dev/null || true
	@sleep 2
	@echo "🚀 重新启动服务..."
	@make -j2 dev-server dev-web

restart-server: ## 重启后端服务
	@echo "🔄 重启后端服务..."
	@pkill -f "go run main.go" 2>/dev/null || true
	@sleep 1
	@make dev-server

restart-web: ## 重启前端服务
	@echo "🔄 重启前端服务..."
	@pkill -f "react-scripts start" 2>/dev/null || true
	@sleep 1
	@make dev-web

server: build-server ## 运行生产模式服务器
	@echo "🚀 启动TwinOS生产服务器..."
	@./$(BUILD_DIR)/$(SERVER_NAME)

web: build-web ## 运行生产模式Web应用(需要静态服务器)
	@echo "📱 Web应用已构建完成，请使用静态文件服务器部署"
	@echo "构建文件位置: $(DIST_DIR)/$(WEB_NAME)"

test: ## 运行所有测试
	@echo "🧪 运行Go测试..."
	@cd $(BACKEND_DIR) && go test -v ./...
	@echo "🧪 运行Node.js测试..."
	@cd $(FRONTEND_DIR) && npm test -- --coverage --watchAll=false

test-server: ## 运行后端测试
	@cd $(BACKEND_DIR) && go test -v ./...

test-web: ## 运行前端测试
	@cd $(FRONTEND_DIR) && npm test -- --coverage --watchAll=false

lint: ## 运行代码检查
	@echo "🔍 检查Go代码..."
	@cd $(BACKEND_DIR) && golangci-lint run 2>/dev/null || golint ./... 2>/dev/null || echo "golint/golangci-lint 未安装"
	@echo "🔍 检查TypeScript代码..."
	@cd $(FRONTEND_DIR) && npm run lint 2>/dev/null || echo "ESLint未配置"

format: ## 格式化代码
	@echo "💅 格式化Go代码..."
	@cd $(BACKEND_DIR) && go fmt ./... && goimports -w .
	@echo "💅 格式化TypeScript代码..."
	@cd $(FRONTEND_DIR) && npm run format 2>/dev/null || echo "Prettier未配置"

install: build ## 安装到系统路径
	@echo "📥 安装TwinOS到系统路径..."
	@sudo cp $(BUILD_DIR)/$(SERVER_NAME) /usr/local/bin/
	@echo "✅ 安装完成，使用 'twin-os-server' 启动服务"

docker-build: ## 构建Docker镜像
	@echo "🐳 构建Docker镜像..."
	@docker build -t $(PROJECT_NAME):$(VERSION) .
	@docker tag $(PROJECT_NAME):$(VERSION) $(PROJECT_NAME):latest
	@echo "✅ Docker镜像构建完成"

docker-run: ## 运行Docker容器
	@echo "🐳 启动Docker容器..."
	@docker run -d \
		--name $(PROJECT_NAME) \
		-p 8080:8080 \
		-v $(PWD)/data:/app/data \
		-v $(PWD)/backups:/app/backups \
		$(PROJECT_NAME):latest

status: ## 查看项目状态
	@echo "📊 TwinOS项目状态"
	@echo "=================="
	@echo "版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Go版本: $(shell go version)"
	@echo "Node.js版本: $(NODE_VERSION)"
	@echo "NPM版本: $(NPM_VERSION)"
	@echo "操作系统: $(GOOS)/$(GOARCH)"
	@echo ""
	@echo "📁 文件结构:"
	@tree -L 2 -I 'node_modules|vendor|*.log' 2>/dev/null || ls -la

backup: ## 创建项目备份
	@echo "💾 创建项目备份..."
	@mkdir -p backups
	@tar -czf backups/twin-os-backup-$(shell date +%Y%m%d-%H%M%S).tar.gz \
		--exclude='node_modules' \
		--exclude='vendor' \
		--exclude='*.log' \
		--exclude='build' \
		--exclude='dist' \
		--exclude='.git' \
		.
	@echo "✅ 备份创建完成"

init: ## 初始化新项目环境
	@echo "🎯 初始化TwinOS开发环境..."
	@make deps
	@make setup-db
	@mkdir -p logs backups
	@cp $(BACKEND_DIR)/.env.example $(BACKEND_DIR)/.env 2>/dev/null || true
	@echo "✅ 开发环境初始化完成"
	@echo ""
	@echo "📝 下一步:"
	@echo "1. 编辑 backend/.env 配置文件"
	@echo "2. 运行 'make dev' 启动开发环境"

release: ## 准备发布版本
	@echo "🚀 准备发布版本..."
	@make clean
	@make test
	@make build
	@mkdir -p $(DIST_DIR)/release
	@tar -czf $(DIST_DIR)/release/twin-os-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz \
		-C $(BUILD_DIR) $(SERVER_NAME) \
		-C ../$(DIST_DIR) $(WEB_NAME) \
		--transform 's,^,twin-os/,' \
		LICENSE README.md
	@echo "✅ 发布版本准备完成: $(DIST_DIR)/release/"

# 开发快捷命令
run: dev ## dev的别名
start: dev ## dev的别名

all: clean deps test build ## 完整的CI流程

version: ## 显示版本信息
	@echo "TwinOS版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Git提交: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
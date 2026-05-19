# LOJ

LOJ 是一个前后端分离的在线评测系统项目，当前仓库已经合并为统一的 monorepo，包含 Vue 前端、Go 后端、Nginx 反向代理以及 Docker Compose 本地编排配置。

## 项目结构

```text
LOJ/
  backend/              # Go + Gin 后端服务
  frontend/             # Vue 前端项目
  nginx/                # Nginx 配置与镜像构建文件
  docker-compose.yml    # 本地开发/部署编排
  GO1_db.sql            # 数据库初始化 SQL
```

## 技术栈

- 前端：Vue 3、Vue Router、Pinia、Element Plus、Naive UI、Axios
- 后端：Go、Gin、GORM、MySQL、Redis、RabbitMQ
- 部署：Docker、Docker Compose、Nginx

## 快速启动

确保本机已安装 Docker 和 Docker Compose，然后在项目根目录执行：

```powershell
docker compose up -d --build
```

服务启动后访问：

```text
http://localhost:8080
```

查看服务状态：

```powershell
docker compose ps
```

停止服务：

```powershell
docker compose down
```

## 前端开发

```powershell
cd frontend
npm install
npm run serve
```

构建前端：

```powershell
npm run build
```

## 后端开发

```powershell
cd backend
go mod download
go run main.go
```

后端主要配置文件位于：

```text
backend/settings.yaml
```

## Git 仓库说明

本项目由原前端仓库和原后端仓库合并而来。合并后统一由根目录 Git 仓库管理，子项目中不再保留独立的 `.git` 或分散的 Git 忽略规则。

建议后续只在项目根目录执行 Git 命令：

```powershell
git status
git add .
git commit -m "your commit message"
git push
```

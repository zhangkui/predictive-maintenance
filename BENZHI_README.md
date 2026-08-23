# PredictiveMaintenance 容器交付说明

`benzhi.Dockerfile` 用于构建 PredictiveMaintenance 的生产镜像。构建过程只复制 Go 模块、服务入口、内部包和日志包，不会把本地数据库、Git 元数据或题库材料写入镜像构建上下文。

## 构建镜像

在项目根目录执行：

```sh
sh build_benzhi_docker.sh
```

默认镜像标签为：

```text
predictive-maintenance:benzhi
```

可通过环境变量覆盖镜像名称和标签：

```sh
IMAGE_NAME=registry.example.com/maintenance IMAGE_TAG=v1 sh build_benzhi_docker.sh
```

## 运行容器

服务使用 SQLite 文件保存运行数据。运行时应挂载一个可写目录到 `/app/data`：

```sh
docker run --rm -p 8080:8080 \
  -e HTTP_ADDR=:8080 \
  -e DB_PATH=/app/data/predictive-maintenance.db \
  -v "$(pwd)/data:/app/data" \
  predictive-maintenance:benzhi
```

容器内默认监听 `8080` 端口。启动后可使用健康检查接口确认服务状态：

```sh
curl --fail http://localhost:8080/healthz
```

预期响应包含 `status` 为 `ok`。

## Docker Compose

现有 `docker-compose.yml` 可直接启动项目：

```sh
docker compose up --build
```

如需使用本交付 Dockerfile 构建 Compose 镜像，可在本地 Compose 覆盖配置中为服务增加：

```yaml
build:
  context: .
  dockerfile: benzhi.Dockerfile
```

# PredictiveMaintenance

工业物联网设备巡检与预测性维护后端 MVP，使用 Go + SQLite。

## 启动

```bash
docker compose up --build
```

默认监听 `http://localhost:8080`，健康检查：`GET /healthz`。

## 示例

```bash
curl http://localhost:8080/api/v1/dashboard/overview
curl http://localhost:8080/api/v1/devices
curl -X POST http://localhost:8080/api/v1/data -H 'Content-Type: application/json' -d '{"deviceId":1,"items":[{"sensorType":"temperature","value":42.5}]}'
```

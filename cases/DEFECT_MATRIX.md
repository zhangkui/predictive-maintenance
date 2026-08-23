# PredictiveMaintenance 30题缺陷矩阵

| ID | 类型 | 功能域 | 单一业务根因 | 生产文件（至少2） | 公开验证重点 | 冲突资源 |
|---|---|---|---|---|---|---|
| BUG-001 | bugfix/其它 | 异常工单 | 持续异常扫描缺少未完成工单幂等约束 | service/ticket_service.go; scheduler/anomaly_scanner.go | 重复扫描只保留一张未完成工单，完成后允许新工单 | ticket create; anomaly scan |
| BUG-002 | bugfix/error | 传感器接收 | 批量写入异常时没有事务性回滚 | service/sensor_service.go; repository/sensor_repo.go | 批次中途失败时传感器和异常记录均不落库 | ingest batch; sensor insert |
| BUG-003 | bugfix/其它 | 阈值检测 | 低于正常下限的数据被当作正常 | detector/threshold.go; service/sensor_service.go | 下限越界、危急下限和正常值分别判定 | threshold detect |
| BUG-004 | bugfix/error | 变化率检测 | 时间间隔无效时变化率错误污染异常判断 | detector/rate.go; util/math.go | 零/负时间间隔安全处理并保持后续样本可判定 | rate calculation |
| BUG-005 | bugfix/其它 | 趋势检测 | 趋势窗口退化导致短期噪声被当成持续恶化 | detector/trend.go; service/prediction_service.go | 多点趋势、噪声和恢复趋势分别验证 | trend window |
| BUG-006 | bugfix/其它 | 组合检测 | 多探头异常只取最后一个结果，丢失最高严重级别 | detector/composite.go; service/anomaly_service.go | 多传感器组合保留最高严重性和原因集合 | composite aggregation |
| BUG-007 | bugfix/其它 | 健康评估 | 缺失传感器维度被按满分计入综合健康度 | service/health_service.go; service/health_history.go | 缺失数据扣分，完整数据权重保持稳定 | health scoring |
| BUG-008 | bugfix/其它 | 健康告警 | 健康度低于阈值只更新分数未创建告警记录 | service/health_service.go; repository/abnormal_repo.go | 首次低分告警、重复评估幂等、恢复状态清理 | health alarm |
| BUG-009 | bugfix/其它 | 寿命预测 | 下降趋势方向被反转，恶化设备得到更长寿命 | service/prediction_service.go; detector/trend.go | 下降/平稳/改善趋势的剩余寿命排序 | prediction trend |
| BUG-010 | bugfix/其它 | 预测置信度 | 样本不足仍给出高置信度预测 | service/prediction_service.go; util/math.go | 样本量、波动度影响置信度且范围受限 | prediction confidence |
| BUG-011 | bugfix/其它 | 维保周期 | 周期计划按创建时间重复生成而不基于最近任务 | service/maintenance_service.go; repository/maintenance_repo.go | 同一计划周期内只生成一次，跨周期可生成 | plan due query |
| BUG-012 | bugfix/error | 维保运行时长 | 运行小时累计查询包含计划区间外记录 | service/runtime_service.go; service/maintenance_service.go | 区间边界、累计阈值和重置后状态 | runtime due |
| BUG-013 | bugfix/其它 | 维保状态 | 任务状态允许跳过合法流转直接完成 | repository/maintenance_repo.go; service/maintenance_service.go | pending→running→completed，非法跳转拒绝 | task transition |
| BUG-014 | bugfix/其它 | 维保日历 | 周末跳过逻辑只检查当前日期不调整下一个工作日 | scheduler/maintenance_creator.go; util/time.go | 周末和连续周末任务顺延到工作日 | schedule date |
| BUG-015 | bugfix/其它 | 设备离线 | 离线检测以最近创建时间而非采集时间判定 | scheduler/offline_detector.go; repository/sensor_repo.go | 延迟写入、旧采集数据和在线设备状态 | offline detection |
| BUG-016 | bugfix/其它 | 设备分组 | 分组删除未处理设备关联，列表出现孤儿设备 | service/group_service.go; repository/device_repo.go | 删除/重建分组后的设备归属一致 | group relation |
| BUG-017 | bugfix/error | 传感器配置 | 配置更新不校验危急范围与正常范围关系 | service/config_service.go; service/validation_service.go | 正常范围、危急范围、阈值关系拒绝非法配置 | config validation |
| BUG-018 | bugfix/其它 | 异常处理 | 异常处理更新未记录处理人和处理时间 | repository/abnormal_repo.go; service/audit_service.go | handle/ignore 均持久化操作者和审计事件 | abnormal handling |
| BUG-019 | bugfix/其它 | 工单关联 | 异常转工单未持久化异常关联 ID | service/ticket_service.go; internal/store/store.go | 工单详情可反查异常，关系唯一 | ticket relation |
| BUG-020 | bugfix/其它 | 时序查询 | 相同采集时间的传感器数据结果不稳定 | repository/sensor_repo.go; service/sensor_service.go | 时间相同按 ID 稳定排序，最新数据可重复读取 | sensor ordering |
| BUG-021 | bugfix/error | 数据保留 | 传感器/健康历史清理跨表删除失败时部分提交 | service/retention_service.go; service/health_history.go | 任一删除失败整体回滚，统计保持一致 | retention transaction |
| BUG-022 | bugfix/其它 | 报表 | 日报统计使用本地时区边界导致跨日数据错归 | service/report_service.go; util/time.go | UTC 边界和设备时区数据归属稳定 | report window |
| BUG-023 | bugfix/其它 | 通知 | 同一异常重复通知没有去重窗口 | service/notification_service.go; service/anomaly_service.go | 重复扫描不重复通知，新异常可通知 | notification dedup |
| BUG-024 | bugfix/其它 | 仪表盘 | 健康排名查询没有稳定排序和分页边界 | service/dashboard_service.go; internal/store/store.go | 分数相同按设备 ID，limit/offset 边界准确 | health ranking |
| BUG-025 | bugfix/error | 设备 API | 更新设备状态绕过合法状态校验 | api/api.go; service/device_service.go | 非法状态拒绝，合法状态和更新时间正常 | device update |
| BUG-026 | bugfix/context | 数据 API | 请求取消后批量接收仍继续写入后续数据 | api/api.go; service/sensor_service.go | context cancel 停止写入并返回取消错误 | request context |
| BUG-027 | bugfix/concurrency | 调度器 | Stop 与任务执行并发时存在重复执行/等待不一致 | scheduler/runner.go; scheduler/anomaly_scanner.go | 停止期间不新增执行，WaitGroup 稳定 | scheduler lifecycle |
| BUG-028 | bugfix/defer | 数据库资源 | 查询迭代错误时 rows 未按统一路径释放/传播 | repository/sensor_repo.go; repository/abnormal_repo.go | 行扫描错误和关闭错误均可观察且资源释放 | rows lifecycle |
| BUG-029 | bugfix/其它 | 设备健康状态 | 传感器异常只记录异常不联动设备故障状态 | service/sensor_service.go; service/device_service.go | 危急异常置故障，恢复后按健康度恢复在线 | device health state |
| BUG-030 | bugfix/其它 | 维护计划触发 | 条件维保触发后未幂等绑定计划与任务 | service/maintenance_service.go; repository/maintenance_repo.go | 健康阈值重复扫描只生成一个任务，不同计划隔离 | condition maintenance |

## 预审结论

- 每题只保留一个根本业务异常。
- 每题至少两个生产文件，目标修复量不少于12行实质生产代码。
- 不使用单数字、单运算符、简单 nil、字段拼写或编译错误凑题。
- BUG-001 当前历史候选仅修改一个文件，不计入本矩阵的最终实现，需重新按本矩阵改造。
- 进入生产注入前必须逐题检查关键函数、SQL写入点、事务动作、缓存键和状态动作冲突。

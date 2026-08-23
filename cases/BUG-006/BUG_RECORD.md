# BUG-006 生产缺陷记录

- case_id: BUG-006
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-006 | bugfix/其它 | 组合检测 | 多探头异常只取最后一个结果，丢失最高严重级别 | detector/composite.go; service/anomaly_service.go | 多传感器组合保留最高严重性和原因集合 | composite aggregation |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug006_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
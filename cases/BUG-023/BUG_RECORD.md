# BUG-023 生产缺陷记录

- case_id: BUG-023
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-023 | bugfix/其它 | 通知 | 同一异常重复通知没有去重窗口 | service/notification_service.go; service/anomaly_service.go | 重复扫描不重复通知，新异常可通知 | notification dedup |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug023_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
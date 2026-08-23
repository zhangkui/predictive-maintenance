# BUG-001 生产缺陷记录

- case_id: BUG-001
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-001 | bugfix/其它 | 异常工单 | 持续异常扫描缺少未完成工单幂等约束 | service/ticket_service.go; scheduler/anomaly_scanner.go | 重复扫描只保留一张未完成工单，完成后允许新工单 | ticket create; anomaly scan |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug001_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
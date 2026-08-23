# BUG-007 生产缺陷记录

- case_id: BUG-007
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-007 | bugfix/其它 | 健康评估 | 缺失传感器维度被按满分计入综合健康度 | service/health_service.go; service/health_history.go | 缺失数据扣分，完整数据权重保持稳定 | health scoring |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug007_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
# BUG-021 生产缺陷记录

- case_id: BUG-021
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-021 | bugfix/error | 数据保留 | 传感器/健康历史清理跨表删除失败时部分提交 | service/retention_service.go; service/health_history.go | 任一删除失败整体回滚，统计保持一致 | retention transaction |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug021_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
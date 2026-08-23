# BUG-008 生产缺陷记录

- case_id: BUG-008
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-008 | bugfix/其它 | 健康告警 | 健康度低于阈值只更新分数未创建告警记录 | service/health_service.go; repository/abnormal_repo.go | 首次低分告警、重复评估幂等、恢复状态清理 | health alarm |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug008_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
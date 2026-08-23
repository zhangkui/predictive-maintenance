# BUG-030 生产缺陷记录

- case_id: BUG-030
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-030 | bugfix/其它 | 维护计划触发 | 条件维保触发后未幂等绑定计划与任务 | service/maintenance_service.go; repository/maintenance_repo.go | 健康阈值重复扫描只生成一个任务，不同计划隔离 | condition maintenance |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug030_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
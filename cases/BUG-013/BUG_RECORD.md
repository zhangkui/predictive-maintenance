# BUG-013 生产缺陷记录

- case_id: BUG-013
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-013 | bugfix/其它 | 维保状态 | 任务状态允许跳过合法流转直接完成 | repository/maintenance_repo.go; service/maintenance_service.go | pending→running→completed，非法跳转拒绝 | task transition |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug013_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
# BUG-011 生产缺陷记录

- case_id: BUG-011
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-011 | bugfix/其它 | 维保周期 | 周期计划按创建时间重复生成而不基于最近任务 | service/maintenance_service.go; repository/maintenance_repo.go | 同一计划周期内只生成一次，跨周期可生成 | plan due query |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug011_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
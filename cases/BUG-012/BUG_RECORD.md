# BUG-012 生产缺陷记录

- case_id: BUG-012
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-012 | bugfix/error | 维保运行时长 | 运行小时累计查询包含计划区间外记录 | service/runtime_service.go; service/maintenance_service.go | 区间边界、累计阈值和重置后状态 | runtime due |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug012_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
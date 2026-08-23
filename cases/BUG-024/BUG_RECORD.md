# BUG-024 生产缺陷记录

- case_id: BUG-024
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-024 | bugfix/其它 | 仪表盘 | 健康排名查询没有稳定排序和分页边界 | service/dashboard_service.go; internal/store/store.go | 分数相同按设备 ID，limit/offset 边界准确 | health ranking |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug024_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
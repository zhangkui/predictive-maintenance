# BUG-022 生产缺陷记录

- case_id: BUG-022
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-022 | bugfix/其它 | 报表 | 日报统计使用本地时区边界导致跨日数据错归 | service/report_service.go; util/time.go | UTC 边界和设备时区数据归属稳定 | report window |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug022_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
# BUG-026 生产缺陷记录

- case_id: BUG-026
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-026 | bugfix/context | 数据 API | 请求取消后批量接收仍继续写入后续数据 | api/api.go; service/sensor_service.go | context cancel 停止写入并返回取消错误 | request context |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug026_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
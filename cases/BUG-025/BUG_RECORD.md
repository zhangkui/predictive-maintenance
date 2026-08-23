# BUG-025 生产缺陷记录

- case_id: BUG-025
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-025 | bugfix/error | 设备 API | 更新设备状态绕过合法状态校验 | api/api.go; service/device_service.go | 非法状态拒绝，合法状态和更新时间正常 | device update |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug025_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
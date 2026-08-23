# BUG-029 生产缺陷记录

- case_id: BUG-029
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-029 | bugfix/其它 | 设备健康状态 | 传感器异常只记录异常不联动设备故障状态 | service/sensor_service.go; service/device_service.go | 危急异常置故障，恢复后按健康度恢复在线 | device health state |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug029_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
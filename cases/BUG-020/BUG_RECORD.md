# BUG-020 生产缺陷记录

- case_id: BUG-020
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-020 | bugfix/其它 | 时序查询 | 相同采集时间的传感器数据结果不稳定 | repository/sensor_repo.go; service/sensor_service.go | 时间相同按 ID 稳定排序，最新数据可重复读取 | sensor ordering |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug020_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
# BUG-015 生产缺陷记录

- case_id: BUG-015
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-015 | bugfix/其它 | 设备离线 | 离线检测以最近创建时间而非采集时间判定 | scheduler/offline_detector.go; repository/sensor_repo.go | 延迟写入、旧采集数据和在线设备状态 | offline detection |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug015_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
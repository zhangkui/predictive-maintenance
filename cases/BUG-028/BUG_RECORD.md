# BUG-028 生产缺陷记录

- case_id: BUG-028
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-028 | bugfix/defer | 数据库资源 | 查询迭代错误时 rows 未按统一路径释放/传播 | repository/sensor_repo.go; repository/abnormal_repo.go | 行扫描错误和关闭错误均可观察且资源释放 | rows lifecycle |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug028_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
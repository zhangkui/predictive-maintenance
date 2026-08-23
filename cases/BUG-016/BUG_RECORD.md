# BUG-016 生产缺陷记录

- case_id: BUG-016
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-016 | bugfix/其它 | 设备分组 | 分组删除未处理设备关联，列表出现孤儿设备 | service/group_service.go; repository/device_repo.go | 删除/重建分组后的设备归属一致 | group relation |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug016_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
# BUG-017 生产缺陷记录

- case_id: BUG-017
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-017 | bugfix/error | 传感器配置 | 配置更新不校验危急范围与正常范围关系 | service/config_service.go; service/validation_service.go | 正常范围、危急范围、阈值关系拒绝非法配置 | config validation |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug017_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
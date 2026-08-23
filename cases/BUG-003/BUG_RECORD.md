# BUG-003 生产缺陷记录

- case_id: BUG-003
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-003 | bugfix/其它 | 阈值检测 | 低于正常下限的数据被当作正常 | detector/threshold.go; service/sensor_service.go | 下限越界、危急下限和正常值分别判定 | threshold detect |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug003_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
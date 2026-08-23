# BUG-004 生产缺陷记录

- case_id: BUG-004
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-004 | bugfix/error | 变化率检测 | 时间间隔无效时变化率错误污染异常判断 | detector/rate.go; util/math.go | 零/负时间间隔安全处理并保持后续样本可判定 | rate calculation |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug004_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
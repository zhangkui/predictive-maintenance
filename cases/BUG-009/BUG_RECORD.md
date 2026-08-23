# BUG-009 生产缺陷记录

- case_id: BUG-009
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-009 | bugfix/其它 | 寿命预测 | 下降趋势方向被反转，恶化设备得到更长寿命 | service/prediction_service.go; detector/trend.go | 下降/平稳/改善趋势的剩余寿命排序 | prediction trend |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug009_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
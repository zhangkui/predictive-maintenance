# BUG-005 生产缺陷记录

- case_id: BUG-005
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-005 | bugfix/其它 | 趋势检测 | 趋势窗口退化导致短期噪声被当成持续恶化 | detector/trend.go; service/prediction_service.go | 多点趋势、噪声和恢复趋势分别验证 | trend window |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug005_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
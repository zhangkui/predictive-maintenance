# BUG-010 生产缺陷记录

- case_id: BUG-010
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-010 | bugfix/其它 | 预测置信度 | 样本不足仍给出高置信度预测 | service/prediction_service.go; util/math.go | 样本量、波动度影响置信度且范围受限 | prediction confidence |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug010_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
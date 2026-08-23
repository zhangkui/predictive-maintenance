# BUG-027 生产缺陷记录

- case_id: BUG-027
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-027 | bugfix/concurrency | 调度器 | Stop 与任务执行并发时存在重复执行/等待不一致 | scheduler/runner.go; scheduler/anomaly_scanner.go | 停止期间不新增执行，WaitGroup 稳定 | scheduler lifecycle |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug027_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
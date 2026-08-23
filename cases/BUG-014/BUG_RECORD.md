# BUG-014 生产缺陷记录

- case_id: BUG-014
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-014 | bugfix/其它 | 维保日历 | 周末跳过逻辑只检查当前日期不调整下一个工作日 | scheduler/maintenance_creator.go; util/time.go | 周末和连续周末任务顺延到工作日 | schedule date |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug014_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
# BUG-019 生产缺陷记录

- case_id: BUG-019
- status: initialized
- scope: 至少修复两个生产 Go 文件，实质生产代码变化不少于 12 行。
- business_design: | BUG-019 | bugfix/其它 | 工单关联 | 异常转工单未持久化异常关联 ID | service/ticket_service.go; internal/store/store.go | 工单详情可反查异常，关系唯一 | ticket relation |
- production_only_initialization: red 与 green 初始化生产代码必须一致；red 仅增加本题回归测试，green 不包含测试。
- verification: TestBug019_BusinessRegression
- workflow: 30 题全部完成前不得启动 Claude。
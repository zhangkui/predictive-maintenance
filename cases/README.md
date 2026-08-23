# Trial Initialization

本目录记录 PredictiveMaintenance 的 30 个真实业务缺陷设计材料。

初始化顺序：

1. `defect` 分支完成全部生产缺陷与唯一公开回归测试。
2. 全量执行 stable-red。
3. 为每题创建 `bugNNN_red` 与 `bugNNN_green`；两者初始化生产代码一致，red 只额外包含本题测试。
4. 两个分支均提交并推送，且通过 `INITIAL_BRANCH_GATE` 后，才能启动 Claude。

当前项目未配置 `origin`，因此禁止宣称已推送或门禁通过。

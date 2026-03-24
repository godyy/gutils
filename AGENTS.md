# AGENTS 指导文件（github.com/godyy/gutils）

面向自动化代理（AI Agents）与开发者的协作指南，确保对本仓库的改动一致、可靠、可维护。

## 项目概览

- 模块：github.com/godyy/gutils
- 语言/版本：Go 1.22
- 现有依赖：github.com/pkg/errors、github.com/rs/xid、golang.org/x/crypto、gopkg.in/check.v1
- 主要包：
  - buffer（bytes、fixed_buffer 等）
  - container（heap、set、skiplist）
  - crypto/password
  - debug（stacktrace）
  - deepcopy
  - flags
  - params（optional）
  - ratelimit（counter、leaky_bucket、token_bucket）
  - worker

## 核心原则

- 简洁稳健：保持 API 简洁、明确，尽量零/低分配。
- 一致性优先：遵循已存在的命名、错误处理与测试模式。
- 并发安全：涉及共享状态的实现要关注竞态与内存可见性。
- 向后兼容：避免破坏性变更；必要时以新增 API 代替修改签名。
- 依赖克制：优先标准库与现有依赖，新增依赖需有充分价值与可维护性。

## 代码风格与约定

- 命名
  - 导出符号使用清晰的驼峰命名；避免不必要缩写。
  - 泛型类型参数使用有意义名称（如 K、V、T、Item）。
- 错误处理
  - 返回 error 为主；保留上下文信息，使用 github.com/pkg/errors 进行包装（Wrap/WithStack）。
  - 不在库中打印日志；将上下文向上传递，由调用方决定记录方式。
- 文档与可读性
  - 为导出类型、函数、接口编写注释，说明语义、边界与复杂度。
  - 在涉及性能权衡的实现中，简述设计意图和限制。
- 性能与内存
  - 避免不必要分配与拷贝；热点路径尽量零分配。
  - 对性能敏感改动提供 Benchmark 与对比数据。
- 并发与竞态
  - 涉及 goroutine/共享状态：确保无泄漏、无死锁，必要时提供 Close/Stop。
  - 使用 `go test -race` 检测数据竞态；时间相关逻辑尽量可注入时间源以降低 flakiness。

## 包级注意事项（示例）

- buffer
  - 关注容量扩展策略与 Endian 读写一致性；遵循 io.Reader/Writer 语义预期。
- container（heap、skiplist、set）
  - 明确时间复杂度；必要时提供示例/基准与对比。
- ratelimit（counter、leaky_bucket、token_bucket）
  - 支持高并发访问；已存在基于 atomic 与 mutex 的实现，保持接口一致与测试覆盖。
- worker
  - 保证可控生命周期与优雅退出；避免 goroutine 泄漏与资源未释放。

## 测试与验证

- 单元测试
  - 运行全部测试：`go test ./...`
  - 覆盖率：`go test -cover ./...`
  - 竞态检测：`go test -race ./...`
- 基准测试
  - 运行：`go test -bench=. ./...`
  - 对性能相关改动提供基准与前后对比，避免回退。
- 测试质量
  - 避免脆弱用例；对时间/随机性使用可控制注入或固定种子。
  - 错误路径与边界条件应有覆盖。

## 依赖与版本

- Go 1.22，模块路径固定为 `github.com/godyy/gutils`。
- 新增依赖需满足：小而稳定、许可证兼容、活跃维护；优先标准库。
- 更新依赖需跑全量测试与竞态检测，确认无行为回退。

## 变更流程

1. 需求澄清：明确 API 设计、复杂度、边界行为与并发模型。
2. 设计对齐：对破坏性风险进行规避（新增 API 而非修改旧签名）。
3. 实现规范：遵循现有风格、错误处理与测试模式。
4. 验证与基准：`go test ./...`、`-race`、必要时 `-bench`。
5. 文档完善：为导出符号与公开行为添加注释与示例。

## 代码审查核对清单

- API 语义清晰、命名一致、无多余暴露。
- 错误处理保留上下文且不吞错。
- 并发安全（无泄漏、无竞态），可通过 `-race`。
- 性能不回退，基准数据充分。
- 单测覆盖边界与错误路径，稳定不脆弱。
- 无不必要的新依赖与反射；跨平台行为一致。

## 常用命令

```bash
go test ./...
go test -race ./...
go test -cover ./...
go test -bench=. ./...
go vet ./...
```

## 安全与合规

- 不引入/泄露任何密钥、凭据或敏感数据。
- 输入校验与边界检查严格，不信任外部输入。
- 保持第三方依赖许可证兼容性与可追溯性。

---

若需要新增通用工具包，请参考现有目录结构与实现风格；提交前确保通过测试与竞态检测，并补充必要文档与基准数据。


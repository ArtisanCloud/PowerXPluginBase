# Research Findings — Plugin Capabilities & Schema Governance

## Decision 1: Capability metadata resides in contracts directory
- **Decision**: Manifest仅保留能力 ID、类型、版本、依赖等最小信息；详细定义（描述、输入/输出、版本历史）集中在 `contracts/capabilities/*.yaml`。
- **Rationale**: 保持 manifest 精简，避免大块 YAML；同时让工具和宿主通过统一目录读取更详细的描述。
- **Alternatives considered**:
  - 在 manifest 中写所有细节 —— 容易膨胀且难以版本化；
  - 把所有信息散在文档 —— 缺乏机器可读性。

## Decision 2: Controlled breaking changes via deprecation windows
- **Decision**: 允许 MAJOR 版本进行破坏性调整，但必须提供适配器/弃用窗口，并在宿主侧用版本 gate 控制升级。
- **Rationale**: 完全禁止破坏性更改不现实；通过文档化流程+宿主 gate 保证安全过渡。
- **Alternatives considered**:
  - 只允许向后兼容变更 —— 会囤积大量历史字段；
  - 完全放任团队自定义 —— 难以做平台级验证。

## Decision 3: Automation targets prescribed as recommendations with requirements on outcome
- **Decision**: 规范提供 `make` 命令参考（如 `make check-capability`, `make check-compat`），并要求 CI 出具校验报告，但团队可采用等效工具，只要输出满足标准。
- **Rationale**: 给出统一入口便于平台集成，同时保留技术选型灵活性。
- **Alternatives considered**:
  - 强制使用特定工具 —— 对不同语言/生态限制较大；
  - 仅写原则无命令 —— 难以推动落地与自动化。

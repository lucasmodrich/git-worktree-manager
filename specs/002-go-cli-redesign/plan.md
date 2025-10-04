
# Implementation Plan: Rename Go CLI Binary to "gwtm"

**Branch**: `002-go-cli-redesign` | **Date**: 2025-10-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/home/modrich/dev/lucasmodrich/git-worktree-manager/specs/002-go-cli-redesign/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
This plan implements a simple but important change: renaming the Go CLI binary from `git-worktree-manager` to `gwtm` for improved user experience and convenience. The shortened name makes the tool easier to type while maintaining all existing functionality. This change affects the binary name, build configuration, documentation, and CI/CD workflows.

## Technical Context
**Language/Version**: Go 1.21+
**Primary Dependencies**: Cobra v1.10.1 (CLI framework)
**Storage**: N/A (CLI tool, no persistent storage)
**Testing**: Go standard testing package with table-driven tests
**Target Platform**: Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64)
**Project Type**: Single (CLI application)
**Performance Goals**: <2 seconds for all commands
**Constraints**: Backward compatibility with existing `git-worktree-manager` binary name during transition
**Scale/Scope**: Simple refactoring affecting build configuration, 5-10 files to update

**User Request**: Rename the Go CLI app binary to "gwtm"

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Multi-Implementation Strategy
- ✅ **PASS**: Binary rename affects only Go implementation. Bash script (`git-worktree-manager.sh`) remains unchanged.
- ✅ **PASS**: Both implementations can coexist; users choose which to use.

### Principle II: Language-Specific Best Practices (Go)
- ✅ **PASS**: Standard Go build practices apply; binary name set via build output flag.
- ✅ **PASS**: No changes to code structure or naming conventions required.

### Principle III: Test-First Development
- ✅ **PASS**: No new functionality added; existing tests validate binary continues to work.
- ✅ **PASS**: Build verification tests will confirm new binary name.

### Principle IV: User Safety & Transparency
- ✅ **PASS**: No destructive operations involved; pure rename.
- ✅ **PASS**: Documentation will clearly communicate the change.

### Principle V: Semantic Release Compliance
- ⚠️ **CONSIDERATION**: Binary name change could be considered a breaking change if users have hardcoded paths.
- ✅ **MITIGATION**: Document as MINOR version bump (new feature: shorter binary name) since old name can be aliased.
- ✅ **PASS**: Conventional commits will be used.

### Principle VI: Backward Compatibility
- ⚠️ **CONSIDERATION**: Users with scripts/aliases using `git-worktree-manager` may break.
- ✅ **MITIGATION**: Documentation will provide migration guide; consider creating symlink or alias in installation.
- ✅ **PASS**: CLI interface (commands, flags) remains identical.

**Initial Constitution Check: PASS** (with documented mitigations for compatibility)

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
/home/modrich/dev/lucasmodrich/git-worktree-manager/
├── cmd/
│   └── git-worktree-manager/     # Main entry point (will build as "gwtm")
│       └── main.go
├── internal/
│   ├── commands/                  # Cobra command implementations
│   ├── config/                    # Configuration management
│   ├── git/                       # Git operations
│   ├── ui/                        # User interface utilities
│   └── version/                   # Version management
├── tests/
│   ├── integration/               # Integration tests
│   └── contract/                  # Contract tests
├── Makefile                       # Build configuration (update binary name)
├── .goreleaser.yml                # Release configuration (update binary name)
├── .github/workflows/
│   ├── test.yml                   # CI test workflow (update binary name)
│   └── release.yml                # Release workflow (update binary name)
└── README.md                      # Documentation (update examples)
```

**Structure Decision**: Single-project Go CLI application. The binary rename affects build configuration files (Makefile, .goreleaser.yml) and CI/CD workflows. Source code structure remains unchanged. The `cmd/git-worktree-manager/` directory name can remain as-is since it's the package path; only the output binary name changes to `gwtm`.

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh claude`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Binary rename is a simple refactoring task, NOT a full feature implementation
- Focus on build configuration, documentation, and CI/CD updates
- No new source code or tests required (existing functionality unchanged)

**Task Categories**:
1. **Build Configuration** (3 tasks):
   - T001: Update Makefile to output `gwtm` binary
   - T002: Update .goreleaser.yml binary name to `gwtm`
   - T003: Update go.mod module path if needed (likely no change)

2. **CI/CD Workflows** (2 tasks):
   - T004: Update `.github/workflows/test.yml` build and test commands
   - T005: Verify `.github/workflows/release.yml` (uses .goreleaser.yml config)

3. **Documentation** (2 tasks):
   - T006: Update README.md with new binary name `gwtm` and migration guide
   - T007: Update CLAUDE.md with binary name reference

4. **Validation** (1 task):
   - T008: Build and smoke test new binary name (./gwtm --help, ./gwtm version)

**Ordering Strategy**:
- Sequential for safety (build config → CI/CD → docs → validation)
- Mark build and doc tasks as [P] for parallel execution (independent files)

**Estimated Output**: 8 numbered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command) ✅
- [x] Phase 1: Design complete (/plan command) ✅
- [x] Phase 2: Task planning complete (/plan command - describe approach only) ✅
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS ✅
- [x] Post-Design Constitution Check: PASS ✅
- [x] All NEEDS CLARIFICATION resolved ✅
- [x] Complexity deviations documented (N/A - simple refactoring) ✅

**Artifacts Generated**:
- [x] research.md (updated with binary rename section 11)
- [x] data-model.md (existing, no changes needed)
- [x] contracts/cli-interface.md (updated with binary name note)
- [x] quickstart.md (updated with binary name note)
- [x] CLAUDE.md (updated via script)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*

<div align="center">

<img src="qa-align-cli/demo.png" alt="QA-Align telemetry engine demo" width="760"/>

# qa-align

**Stop copying test names into Jira by hand.**  
`qa-align` scans your codebase, reads your docstring annotations, counts Git churn, and exports a risk-ranked `telemetry.json` audit trail — in under a second, with zero runtime dependencies.

[![Go Version](https://img.shields.io/badge/go-1.21%2B-00ADD8?logo=go)](https://go.dev)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-lightgrey)]()

</div>

---

## The Problem

Every sprint, someone on the QA team manually copies test method names, requirement IDs, and risk levels into a spreadsheet or Jira ticket. That person is:

- ❌ Doing it wrong — test code drifted since the last copy-paste  
- ❌ Doing it slow — 2 hours per release cycle, minimum  
- ❌ Not doing security traceability — a CISO audit finds gaps instantly

`qa-align` eliminates that entirely. It reads the truth directly from your source files.

---

## 60-Second Quickstart

### 1 — Install (from source, ~10 seconds)

```bash
git clone https://github.com/JRE503/TestGate-cli
cd qa-align-cli
go build -o qa-align
```

> **Requires:** [Go 1.21+](https://go.dev/dl). No Docker, no Node, no Python.

### 2 — Annotate one test (takes 30 seconds)

Add three comment lines above any test function in Python, TypeScript, JavaScript, Java, Kotlin, or C++:

```python
# [What]: Verifies that a locked account cannot authenticate
# [Why]: Account lockout policy must be enforced after N failed attempts.
# [Reference]: SEC-110
def test_locked_account_cannot_authenticate():
    ...
```

TypeScript / JS also supported via `@test`, `@description`, `@issue`:

```typescript
// @test Rejects duplicate email registration
// @description System must enforce uniqueness before DB write
// @issue USR-302
function test_create_user_duplicate_email_rejected() { ... }
```

### 3 — Run

```bash
./qa-align --dir ./your-project
```

### 4 — Read the output

```json
[
  {
    "test_method": "test_locked_account_cannot_authenticate",
    "file_path": "tests/auth/test_authentication.py",
    "what": "Verifies that a locked account cannot authenticate regardless of credentials",
    "why": "Account lockout policy must be enforced after N failed attempts.",
    "requirement_id": "SEC-110",
    "change_frequency_30_days": 7,
    "calculated_risk_score": 6
  }
]
```

Ship `telemetry.json` straight into your CI artifact store, your Jira importer, or your CISO's inbox.

---

## Annotation Format Reference

| Language | `[What]` / `@test` | `[Why]` / `@description` | `[Reference]` / `@issue` |
|---|---|---|---|
| Python | `# [What]: ...` | `# [Why]: ...` | `# [Reference]: TICKET-123` |
| TypeScript / JS | `// @test ...` | `// @description ...` | `// @issue TICKET-123` |
| Java / Kotlin | `// @test ...` | `// @description ...` | `// @issue TICKET-123` |
| C++ | `// @test ...` | `// @description ...` | `// @issue TICKET-123` |

---

## Risk Score Formula

```
Risk = Impact × Frequency × Probability
```

| Factor | Low (1) | Medium (2) | High (3) |
|---|---|---|---|
| **Impact** | `utils/`, `styles/` paths | All others | `auth`, `billing`, `security` paths |
| **Frequency** | churn ≤ 2 commits/30d | churn 3–10 | churn > 10 |
| **Probability** | — | 2 (fixed baseline) | — |

Max score = **18**. Anything ≥ 12 should be in your regression suite every release.

---

## CI Integration

Drop this into `.github/workflows/qa-telemetry.yml`:

```yaml
name: QA Telemetry

on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Build qa-align
        run: go build -o qa-align
      - name: Run telemetry scan
        run: ./qa-align --dir .
      - name: Upload telemetry artifact
        uses: actions/upload-artifact@v4
        with:
          name: telemetry
          path: telemetry.json
```

---

## Project Structure

```
qa-align-cli/
├── main.go                        # Entry point
├── cmd/root.go                    # CLI command router (cobra)
├── internal/
│   ├── parser/python_parser.go    # Polyglot annotation scanner
│   ├── gitops/churn.go            # Git log churn calculator
│   ├── risk/matrix.go             # Risk = Impact × Frequency × Probability
│   └── schema/normalizer.go      # TestCaseMetadata type + JSON serializer
├── qa-align.json                  # Local config profile
└── demo.svg                       # Animated terminal demo
```

---

| Extension | Language | Frameworks |
|---|---|---|
| `.py` | Python | pytest, unittest |
| `.ts` | TypeScript | Jest, Vitest |
| `.js` | JavaScript | Jest, Mocha |
| `.java` | Java | JUnit, TestNG |
| `.kt` | Kotlin | JUnit5, Kotest |
| `.cpp` | C++ | Google Test, CppUTest |
| `.c` | C | **ESP-IDF (Unity)**, Zephyr (ztest), bare metal |

---

## Embedded / ESP-IDF Support

`qa-align` natively parses ESP-IDF **Unity** `TEST_CASE` macros. Add annotations above the macro:

```c
// @test Verifies NVS read returns correct value after write
// @description Core data persistence. Failure means device loses config on reboot.
// @issue FW-101
TEST_CASE("nvs_read_after_write", "[nvs]")
{
    // ...
}
```

The test name is extracted directly from the first string argument of `TEST_CASE`. Works on any `.c` file in your ESP-IDF project tree.

### Embedded Framework Support Matrix

| Framework | Status | Detection pattern |
|---|---|---|
| ESP-IDF Unity | ✅ | `TEST_CASE("name", "[tag]")` |
| Zephyr ztest | ✅ | `void test_name()` style |
| CppUTest | ✅ | `.cpp` + `void test_name()` |
| Google Test | ✅ | `.cpp` + `void test_name()` |
| Arduino bare metal | ✅ | `.c/.cpp` + `// @test` annotations |


## License

Apache 2.0 © 2025 JRE503 — Built with Go, `cobra`, and `go-git`.

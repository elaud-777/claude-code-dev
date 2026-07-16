# 여러 PC/OS에서 같은 작업 이어가기 — Claude Code on the web

여러 대의 PC(Windows/Mac/Linux 혼용)를 오가며 **같은 작업을 이어서** 하고 싶다면, [SYNC.md](SYNC.md)의 로컬 상태 파일 동기화(Syncthing)보다 **Claude Code on the web**을 쓰는 게 정확히 이 목적에 맞습니다. 실행 자체가 클라우드 VM 한 곳에서만 일어나므로, 여러 기기의 로컬 상태를 동기화할 필요가 아예 없습니다.

## 왜 이게 SYNC.md보다 나은가

| | Remote Control | Claude Code on the web | SYNC.md (Syncthing) |
|---|---|---|---|
| 실행 위치 | 내 로컬 PC (계속 켜져 있어야 함) | Anthropic 클라우드 VM | 각 PC (상태를 서로 복제) |
| 여러 PC에서 접속 | 호스트 PC 하나에 의존 | 어느 기기든 브라우저만 있으면 동일 | 동시 실행 금지 규칙 필요 |
| 크로스 OS | 상관없음 (원격 조종일 뿐) | 상관없음 | ⚠️ 셸 스냅샷 등 OS별 가정 때문에 비권장 |
| 요구사항 | Pro/Max/Team/Enterprise 구독 | Pro/Max/Team 구독 (Enterprise는 프리미엄 seat) | 없음(무료) |
| 로컬 파일 접근 | 가능 | 불가 (GitHub 저장소 기준만) | 가능 |

**Remote Control**은 "이미 로컬에서 실행 중인 내 세션을 다른 기기에서 원격 조종"하는 기능이라 호스트 PC가 항상 켜져 있어야 합니다 — "하나의 원격 환경"이 아니라 "내 PC 원격 조종"에 가깝습니다. 반대로 **Claude Code on the web**은 실행 환경 자체가 클라우드에 있으므로 이 프로젝트처럼 여러 PC를 오가며 작업하기에 적합합니다.

## 설정 절차

1. https://claude.ai/code 에서 GitHub 계정 연동
2. 이 저장소(`elaud-777/claude-code-dev`)를 연결
3. 웹에서 새 세션 시작 — 세션마다 새 브랜치로 작업이 진행됨
4. 다른 PC/OS에서도 동일하게 https://claude.ai/code 접속 → 같은 저장소의 세션 목록에서 이어서 진행

## 로컬 ↔ 웹 전환

- `--teleport`: 웹 세션의 작업을 로컬로 가져와서 이어서 진행
- `--cloud`: 로컬에서 하던 작업을 웹 세션으로 올림

## 한계 — 왜 이 저장소의 CLAUDE.md/OpenSpec 구조가 중요한지

웹 세션은 **로컬 파일시스템에 접근할 수 없고, GitHub 저장소에 커밋된 내용만** 봅니다. 즉:
- 로컬 전용 MCP 서버, 로컬 설정(`claude-cli-setup/settings.template.json` 등)은 웹 세션에 적용되지 않음
- 하지만 이 프로젝트는 이미 규칙(`CLAUDE.md`)과 결정 이력(`openspec/specs/`, `openspec/changes/archive/`)을 **저장소에 커밋**해뒀으므로, 웹 세션에서 새로 시작해도 이 문서들을 읽고 같은 맥락으로 작업을 이어갈 수 있습니다. 이게 바로 이 구조가 "여러 PC에서 같은 작업 이어가기"에 잘 맞는 이유입니다.

## 언제 SYNC.md(Syncthing) 방식을 쓰나

- **같은 OS끼리**(예: Windows PC 두 대) 오가며, 대화 이력/세션 원본 자체를 그대로 보존하고 싶을 때
- 로컬 파일시스템 접근이 꼭 필요한 작업(로컬 전용 MCP, 로컬 DB 등)이 있을 때
- 구독 등급이 Remote Control/웹 세션 요구사항에 못 미칠 때

크로스 OS 조합이라면 이 문서(Claude Code on the web)를 우선 권장합니다.

# 두 PC 간 Claude Code 런타임 상태 지속 동기화

> **여러 OS(Windows/Mac/Linux)를 섞어서 쓴다면 이 문서보다 [REMOTE-WORK.md](REMOTE-WORK.md)(Claude Code on the web)를 먼저 고려하세요.** 이 문서(Syncthing 방식)는 셸 스냅샷 등 OS별 가정이 있어 크로스 OS 조합에는 비권장이며, **같은 OS끼리**만 쓰는 걸 전제로 합니다.

`install.ps1`/`install.sh`가 처리하는 "환경 설정 재현"(모델/플러그인 목록/테마 등)에 더해, **대화 이력·세션 캐시·사용 통계 같은 런타임 상태**까지 두 PC가 계속 같은 상태를 유지하도록 파일 동기화 도구([Syncthing](https://syncthing.net), 무료/오픈소스)로 설정하는 방법입니다.

## 전제 조건 (반드시 지킬 것)

> **두 PC에서 Claude Code를 동시에 실행하지 않는다.** 한쪽 PC 작업을 끝내고, Syncthing이 양쪽 다 "최신 상태(Up to Date)"로 표시되는 걸 확인한 뒤에만 다른 PC에서 Claude Code를 켠다.

이 규칙을 지키면 아래 위험이 사실상 사라집니다:
- `~/.claude.json`(카운터/캐시/프로젝트 목록이 뒤섞인 단일 JSON)을 두 프로세스가 동시에 써서 깨지는 문제
- OAuth 토큰 갱신이 겹쳐서 양쪽 다 로그아웃되는 문제

지키지 않으면(두 PC를 동시에 켜두면) Syncthing은 파일이 깨지는 대신 `history.sync-conflict-<날짜>.jsonl` 같은 **충돌 사본을 남기고 원본은 보존**합니다(무손실이지만 수동 병합이 필요해짐).

## 동기화 대상 정리

| 범위 | 동기화 여부 | 이유 |
|---|---|---|
| `history.jsonl`, `projects/`, `sessions/`, `session-env/`, `file-history/`, `jobs/`, `tasks/`, `plans/` | ✅ 동기화 | 대화 이력/세션 데이터. 대부분 세션 UUID/PID별로 폴더가 나뉘어 있어 두 PC가 만든 파일끼리 충돌할 일이 거의 없음 |
| `.claude.json`, `.claude.json.backup`, `settings.json`, `.credentials.json` | ✅ 동기화 | 설정/통계/로그인 토큰. **동시 실행 금지 규칙이 특히 중요한 대상** |
| `backups/` | ✅ 동기화 | `.claude.json` 자동 백업, 타임스탬프 파일명이라 충돌 없음 |
| `plugins/` 전체 | ❌ 동기화 안 함 | `marketplaces/`는 GitHub 클론본, `cache/`는 그 파생 캐시. 절대경로(`installLocation`)가 사용자명에 박혀 있어 다른 PC(사용자명이 다르면)에서 깨짐. **대신 `install.ps1`/`install.sh`가 이미 `settings.json`의 `enabledPlugins`로 자동 재설치를 처리**하므로 파일 동기화가 불필요 |
| `daemon/`, `shell-snapshots/` | ❌ 동기화 안 함 | 지금 실행 중인 프로세스 PID/제어 키에 종속. 다른 PC로 옮겨도 의미 없고, 매 세션 자동 재생성됨 |
| `cache/` (changelog, closed-issues 등) | ❌ 동기화 안 함 | 순수 네트워크 재조회 캐시 |
| 로그 파일(`daemon.log`), `.last-cleanup`, `.last-update-result.json` | ❌ 동기화 안 함 | 순수 로그/타임스탬프 마커 |

## 설정 방법 (Syncthing)

1. 양쪽 PC에 [Syncthing](https://syncthing.net/downloads/) 설치, 서로를 원격 기기로 추가(Device ID 교환)
2. 새 폴더 공유 추가 시 **폴더 경로를 홈 디렉토리**로 지정 (Windows: `C:\Users\<사용자명>`) — `.claude/`와 `.claude.json`이 서로 다른 위치(홈 디렉토리 바로 아래 vs 그 하위)에 있어서, 홈 디렉토리를 공유 루트로 잡고 `.stignore`로 필요한 것만 골라내는 방식이 가장 깔끔합니다
3. 해당 폴더의 "Ignore Patterns" 설정에 이 폴더의 [`.stignore`](.stignore) 내용을 붙여넣기 (또는 홈 디렉토리에 `.stignore` 파일로 직접 저장)
4. Syncthing UI에서 "미리보기(Preview)"로 실제 동기화 대상 파일 목록을 한 번 확인 — 패턴 문법은 버전마다 미묘하게 다를 수 있어 반드시 육안 확인 권장
5. 두 PC 모두 "Up to Date"가 될 때까지 기다렸다가 작업 전환

## 안전장치: 전환 전 체크 스크립트

`check-before-switch.ps1`/`.sh`를 실행하면 현재 PC에서 Claude Code(또는 그 daemon)가 아직 실행 중인지 확인해줍니다. 다른 PC로 넘어가기 전에 실행하세요.

## 사용자명이 다른 두 PC라면

위 표의 `plugins/` 제외 방침 덕분에 사용자명이 달라도 문제없습니다. 다만 `.claude.json` 안에도 절대경로를 참조하는 필드(예: 프로젝트 경로 목록)가 있을 수 있어, 두 PC의 프로젝트 저장 위치(`D:\claude-code-dev` 등)를 동일하게 맞추는 걸 권장합니다.

# Claude Code CLI 환경 재현 스크립트

이 폴더는 TaskFlow 앱과 무관합니다 — 이 PC에 구성된 **Claude Code CLI 자체의 환경**(설치 방식, 모델/언어/이펙트 레벨 설정, 사용 중인 플러그인 10개)을 다른 PC에서도 동일하게 만들기 위한 스크립트입니다.

## 포함된 것 / 안 된 것

**포함:**
- Claude Code CLI 네이티브 설치 (`install.ps1`/`install.sh`)
- 전역 `settings.json` 템플릿(`settings.template.json`): 모델(`sonnet`), 언어(한국어), `effortLevel: low`, 다크 테마, `askUserQuestionTimeout`, `enabledPlugins`(claude-hud, code-review, context7, feature-dev, frontend-design, playwright, skill-creator, superpowers, neon, vercel), `extraKnownMarketplaces`(claude-hud, claude-plugins-official)

**`install.ps1`/`install.sh`에는 절대 포함 안 함 (보안):**
- API 키 (`apiKeyHelper`) — 원본 `~/.claude/settings.json`에 평문으로 들어있던 필드라 템플릿에서 의도적으로 제거했습니다
- `~/.claude/.credentials.json` (OAuth 토큰)
- `~/.claude.json`의 `machineID`/`userID`/`oauthAccount` 등 머신별 상태

이 파일들은 이 설치 스크립트로 옮길 대상이 아니라, 새 PC에서 `claude` 실행 후 로그인하면 자동으로 새로 생성되는 값입니다. (대화 이력/세션/통계/로그인 토큰까지 두 PC 간 지속적으로 동기화하고 싶다면 아래 [SYNC.md](SYNC.md)를 별도로 참고하세요 — 여긴 보안·충돌 위험이 있어 명시적으로 분리해뒀습니다.)

## 사용법

### Windows

```powershell
.\install.ps1
```

### macOS/Linux

```bash
chmod +x install.sh
./install.sh
```

두 스크립트 모두:
1. `claude` CLI가 없으면 공식 설치 스크립트(`https://claude.sh`)로 설치
2. `~/.claude/settings.json`이 없으면 `settings.template.json`을 그대로 복사, 있으면 **덮어쓰지 않고** 백업만 남긴 뒤 수동 병합을 안내 (기존 PC별 설정을 실수로 지우지 않기 위함)

## 설치 후 수동으로 해야 하는 것

1. `claude` 실행 후 로그인 (또는 API 키를 쓴다면 본인 `ANTHROPIC_API_KEY`를 직접 설정)
2. 첫 실행 시 `settings.json`의 `enabledPlugins`/`extraKnownMarketplaces`를 보고 Claude Code가 마켓플레이스와 플러그인을 자동으로 받아옵니다
3. 프로젝트별 메모리(`~/.claude/projects/<project>/memory/`)는 이 스크립트 대상이 아닙니다 — 프로젝트 규칙은 각 저장소의 `CLAUDE.md`에 커밋해두는 걸 권장합니다(이 저장소의 루트 `CLAUDE.md` 참고)

## 여러 PC/OS에서 같은 작업을 이어가려면

- **여러 OS를 섞어 쓴다면(Windows ↔ Mac/Linux 등)**: [REMOTE-WORK.md](REMOTE-WORK.md) — Claude Code on the web으로 클라우드에서 실행하고 어느 기기든 브라우저로 접속. 로컬 상태 동기화가 아예 필요 없어 가장 안전합니다.
- **로컬 PC에서 계속 돌리던 세션을 폰/다른 기기에서 잠깐 들여다보거나 지시만 하고 싶다면**: [REMOTE-CONTROL.md](REMOTE-CONTROL.md) — Remote Control 상세 가이드(활성화 절차, 보안 모델, 권한 처리, 한계점)
- **같은 OS끼리, 대화 이력/세션 원본 자체를 그대로 보존하고 싶다면**: [SYNC.md](SYNC.md) — Syncthing 기반 지속 동기화. 동시 실행 금지 규칙과 정확한 include/exclude 목록(`.stignore`)을 포함합니다.

## 이 스크립트가 재현하지 못하는 것 (SYNC.md 미설정 시)

- 대화 이력, 세션 캐시, 플러그인 사용 통계 등 순수 런타임 상태
- 이 PC에서 이미 설치된 플러그인 버전(예: superpowers 6.1.1) — 새 PC에서는 마켓플레이스 최신 버전이 설치됩니다. 버전을 고정하고 싶다면 별도로 확인이 필요합니다

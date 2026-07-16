## Why

이 저장소에는 아직 Go 개발을 위한 기본 환경이 갖춰져 있지 않다. 앞으로 진행할 Go 기반 작업(구체 기능은 추후 결정)을 시작하기 전에, 모듈 초기화, 코드 포맷팅/린트, 테스트, 빌드 자동화 등 표준적인 개발 기반을 먼저 마련해야 한다.

## What Changes

- Go 모듈(`go.mod`) 초기화 및 기본 디렉터리 구조(`cmd/`, `internal/` 등) 구성
- 코드 포맷팅 및 린트 도구 설정 (`gofmt`, `golangci-lint` 등)
- 테스트 실행 규칙 및 기본 테스트 스캐폴딩 추가
- 로컬 빌드/실행/테스트를 위한 Makefile 또는 스크립트 추가
- CI 워크플로우(빌드·린트·테스트 자동 실행) 추가

## Capabilities

### New Capabilities
- `go-project-setup`: Go 모듈 초기화, 디렉터리 구조, 빌드/린트/테스트 도구 체인 및 CI 파이프라인을 포함한 기본 개발 환경 구성

### Modified Capabilities
(없음 — 새 저장소이므로 기존 스펙 없음)

## Impact

- 영향받는 코드: 저장소 루트 (신규 `go.mod`, `Makefile`, 린트 설정 파일, CI 워크플로우 파일 추가)
- 영향받는 의존성: Go 툴체인, `golangci-lint`, CI 플랫폼(예: GitHub Actions)
- 구체적인 기능/서비스 범위는 이후 별도 변경(change)에서 결정

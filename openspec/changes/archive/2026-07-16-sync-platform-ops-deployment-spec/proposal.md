## Why

`platform-ops` 스펙은 배포 파이프라인(Vercel)과 환경 분리(SQLite ↔ Neon Postgres)를 여전히 원래 계획대로 명시하고 있다. 그러나 실제로는 Docker Compose(로컬 Postgres 컨테이너 기본값) + Kubernetes/Helm 배포로 완전히 대체되었고, Vercel은 한 번도 실제로 쓰인 적이 없다. 또한 Tailwind CSS와 Swagger UI를 로컬 자산으로 벤더링해 외부 CDN 의존성을 제거한 결정도 스펙에 전혀 반영되어 있지 않다. 코드/인프라 작업이 OpenSpec 절차 밖에서(도구성 작업으로 간주되어) 진행되면서 스펙과 실제 시스템이 어긋난 상태이며, 이번 변경으로 스펙을 실제 구현과 동기화한다.

## What Changes

- `배포 파이프라인` 요구사항을 Vercel 기준에서 Docker Compose(로컬 개발) + Kubernetes/Helm(클러스터 배포) 기준으로 교체
- `환경 분리` 요구사항을 "SQLite ↔ Neon Postgres" 대신 "SQLite 또는 로컬/외부 Postgres, `DATABASE_URL` 스킴으로 자동 분기"로 갱신하고, 로컬 Docker Postgres 컨테이너가 기본값임을 명시
- 신규 요구사항 추가: 프론트엔드/API 문서 정적 자산은 외부 CDN 없이 로컬에 번들링되어야 한다(Tailwind CSS 로컬 빌드, Swagger UI `go:embed` 내장)
- 코드/설정 변경은 없음 — 이번 변경은 이미 구현된 상태를 스펙에 반영하는 문서화 작업

## Capabilities

### New Capabilities
(없음)

### Modified Capabilities
- `platform-ops`: 배포 파이프라인(Vercel → Docker/K8s/Helm), 환경 분리(Neon 필수 → 로컬 Postgres 컨테이너 기본값, 외부 DB는 선택), 정적 자산의 외부 CDN 미사용 요구사항 추가

## Impact

- 영향받는 문서: `openspec/specs/platform-ops/spec.md`만 갱신 (코드/인프라는 이미 해당 상태로 구현 완료됨: `docker-compose.yml`, `k8s/`, `helm/taskflow/`, `frontend/tailwind.css`, `backend/cmd/server/swaggerui/`)
- 영향받지 않는 것: `auth`, `team-management`, `kanban`, `team-chat` capability — 동작 변경 없음

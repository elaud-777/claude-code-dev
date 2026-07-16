## Context

`migrate-backend-to-go` 아카이브 이후, 백엔드/프론트엔드 인프라 작업(Docker Compose, Kubernetes 매니페스트, Helm 차트, Tailwind/Swagger UI 로컬 벤더링, 로컬 Postgres 컨테이너 기본값 전환)이 별도 OpenSpec change 없이 직접 진행되었다. 그 결과 `openspec/specs/platform-ops/spec.md`가 실제 배포 방식(Vercel)과 DB 정책(Neon 필수)을 더 이상 반영하지 않는다. 이번 변경은 순수 문서 동기화이며, 코드나 인프라 설정은 이미 목표 상태로 구현되어 있다.

## Goals / Non-Goals

**Goals:**
- `platform-ops` 스펙의 `배포 파이프라인`, `환경 분리` 요구사항을 실제 구현(Docker Compose + K8s/Helm, 로컬 Postgres 기본값)과 일치시킨다
- 정적 자산 로컬 번들링(외부 CDN 미사용) 결정을 신규 요구사항으로 명문화한다

**Non-Goals:**
- 코드, Dockerfile, k8s/Helm 매니페스트, 프론트엔드 자산 변경 — 이미 완료된 상태
- 새로운 기능이나 배포 방식 도입

## Decisions

### 스펙만 갱신하고 재구현은 하지 않음
이미 검증된(go test, docker compose config, helm lint/template, 브라우저 스크린샷) 구현이 있으므로, 이번 변경의 tasks는 문서 편집만 포함한다. 코드 변경 작업 항목을 만들지 않는다.

### `배포 파이프라인` 요구사항 재정의
Vercel Serverless Functions 대신 "Docker Compose(로컬) + Kubernetes/Helm(클러스터)" 두 경로를 모두 요구사항으로 인정한다. 두 경로 모두 이미지 빌드와 기동 방식이 다르므로 시나리오를 각각 명시한다.

### `환경 분리` 요구사항 재정의
"로컬 SQLite ↔ 운영 Neon Postgres" 이분법을 "SQLite 또는 Postgres(로컬 컨테이너가 기본값, 외부 관리형 Postgres는 선택)"로 넓힌다. `DATABASE_URL` 스킴 기반 자동 분기 로직(코드 변경 없음)은 그대로 유지된다.

## Risks / Trade-offs

- [Risk] 향후 인프라 작업이 다시 OpenSpec 절차 밖에서 진행되면 동일한 드리프트가 재발할 수 있음 → Mitigation: 사용자가 직접 "오픈스펙 내용도 바꾸었나?"로 드리프트를 발견한 사례를 계기로, 인프라성 변경도 규모가 있으면 propose 단계를 거치는 것을 권장 (다만 이번 변경 자체가 그 정책을 강제하지는 않음)

## Why

현재 TaskFlow MVP 백엔드는 Python(FastAPI)으로 구현되어 있다. 이를 Go 기반으로 다시 구현하여 Go 스택 학습/실무 적용을 진행하고자 한다. 프론트엔드와 API 사용자 입장에서는 동작이 동일해야 하므로, 기존에 확정된 API 계약(엔드포인트 18개, DB 스키마 4테이블, 에러 응답 표준)은 그대로 유지하고 구현 언어/런타임만 교체한다.

## What Changes

- **BREAKING** (구현 내부 기준): 기존 Python(FastAPI + SQLAlchemy) 백엔드 코드(`backend/`)를 제거하고 Go 백엔드로 완전히 교체한다. API를 호출하는 클라이언트(프론트엔드) 입장에서는 엔드포인트/요청·응답 형식이 동일하므로 호출 방식에는 영향이 없다.
- Go 웹 프레임워크로 `chi` 라우터를 사용하고, DB 접근은 `sqlc` 또는 GORM 중 구현 단계에서 확정, JWT는 `golang-jwt/jwt`, 비밀번호 해시는 `golang.org/x/crypto/bcrypt`를 사용한다
- 기존 18개 API 엔드포인트(Auth 4 + Team 5 + Task 6 + Chat 3)를 Go로 동일하게 재구현한다
- 기존 DB 스키마 4테이블(users/teams/tasks/messages)과 인덱스를 동일하게 유지한다 (로컬 SQLite / 운영 Neon Postgres)
- 에러 응답 표준 `{ error: { code, message, meta? } }`을 동일하게 유지한다
- API 문서화 방식을 FastAPI 자동 OpenAPI에서 `swaggo/swag` 기반 어노테이션 방식으로 전환한다 (Swagger UI 제공은 동일하게 유지)
- 테스트 프레임워크를 pytest에서 Go 표준 `testing` 패키지(+ `net/http/httptest`)로 전환하며, 기존 핵심 플로우(회원가입/로그인/칸반 CRUD/채팅) 테스트 커버리지를 동일하게 이식한다
- 프론트엔드(`frontend/`)는 API 계약이 동일하므로 수정하지 않는다

## Capabilities

### New Capabilities
(없음 — 기능 요구사항은 변경되지 않으며, 기존 capability의 구현 세부사항만 변경)

### Modified Capabilities
- `platform-ops`: 기술 스택(백엔드 언어/프레임워크/DB 라이브러리), API 문서화 도구(FastAPI 자동 OpenAPI → swaggo), 테스트 프레임워크(pytest → Go testing)가 변경된다. 에러 응답 표준, 환경 분리(DATABASE_URL), 배포 파이프라인 요구사항 자체는 변경되지 않는다.

`auth`, `team-management`, `kanban`, `team-chat` capability는 요구사항(동작)이 변경되지 않으므로 델타 스펙을 작성하지 않는다 — 구현 언어만 바뀔 뿐 API 계약과 시나리오는 `add-taskflow-mvp`에서 정의한 그대로 유지된다.

## Impact

- 영향받는 코드: `D:\claude-code-dev\backend` 디렉토리 전체 교체 (Python 코드/venv/requirements.txt 제거, Go 모듈/코드로 대체)
- 영향받지 않는 코드: `D:\claude-code-dev\frontend` (API 계약 동일, 수정 불필요)
- 신규 의존성: Go 툴체인, chi, sqlc/GORM(택1), golang-jwt, bcrypt(x/crypto), swaggo/swag, sqlite3/postgres 드라이버
- 제거되는 의존성: FastAPI, SQLAlchemy, pydantic, python-jose, passlib, pytest 등 Python 의존성 일체
- 배포 영향: 로컬(SQLite)/운영(Neon Postgres) 환경 분리 방식은 동일하게 유지되나, Vercel 배포 시 Python Serverless Functions 대신 Go 바이너리(또는 Go 런타임) 기반 배포로 전환 필요

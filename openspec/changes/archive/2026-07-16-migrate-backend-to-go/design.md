## Context

`add-taskflow-mvp` change에서 Python(FastAPI + SQLAlchemy) 백엔드로 TaskFlow MVP를 구현했다(18개 API, DB 4테이블, pytest 18개 통과, 프론트엔드와 연동 검증 완료. 배포(Neon/Vercel)만 보류 상태). 이번 변경은 그 구현을 Go로 다시 작성하는 것이며, API 계약(엔드포인트 경로/메서드, 요청·응답 스키마, 에러 코드)은 `auth`/`team-management`/`kanban`/`team-chat` 스펙에 정의된 그대로 유지한다. 프론트엔드는 수정하지 않으므로 Go 백엔드는 기존 프론트엔드가 기대하는 응답 형식을 정확히 재현해야 한다.

## Goals / Non-Goals

**Goals:**
- 기존 18개 API를 Go로 재구현하되 요청/응답 바디, 상태 코드, 에러 코드가 Python 버전과 100% 동일하게 동작
- 기존 프론트엔드(`frontend/`) 코드 변경 없이 Go 백엔드와 그대로 연동
- 로컬 SQLite / 운영 Neon Postgres 환경 분리를 `DATABASE_URL`로 유지
- Swagger UI를 통한 API 문서화 및 JWT 인증 테스트 지원 유지 (도구만 swaggo로 전환)
- 핵심 플로우(회원가입/로그인/칸반 CRUD/채팅) Go 테스트로 이식

**Non-Goals:**
- 새로운 기능 추가나 API 계약 변경 (이번 변경 범위 아님)
- 프론트엔드 수정
- 배포(Vercel/Neon) 자체는 `add-taskflow-mvp`에서 이미 보류된 항목이며 이번 변경에서도 별도로 다루지 않음(구현 완료 후 필요 시 후속 처리)

## Decisions

### 1. 웹 프레임워크: `chi`
표준 라이브러리 `net/http`와 호환되는 경량 라우터인 `chi`를 사용한다. 미들웨어(JWT 인증, CORS, 로깅)를 체이닝하기 쉽고 러닝 커브가 낮다.
- 대안 검토: `gin`(더 무겁고 자체 컨텍스트 타입 사용), 표준 `net/http`만 사용(라우팅 패턴 매칭이 번거로움) — MVP 규모에는 chi가 적당한 절충안.

### 2. DB 접근: `database/sql` + `sqlx`
쿼리를 직접 작성하고 구조체에 매핑하는 `sqlx`를 사용한다. SQLite(로컬)와 Postgres(운영) 방언 차이는 쿼리 레벨에서 최소화(표준 SQL 위주)하고, 드라이버만 `DATABASE_URL` 스킴에 따라 전환한다.
- 대안 검토: GORM(러닝커브는 낮지만 마법 같은 동작으로 쿼리 예측이 어려움), sqlc(컴파일 타임 안전성은 좋으나 코드 생성 파이프라인 추가 필요) — MVP 규모와 팀의 Go 숙련도를 고려해 sqlx로 단순하게 시작하고, 필요시 추후 sqlc로 전환 가능하도록 SQL을 파일로 분리해둔다.

### 3. 인증: `golang-jwt/jwt/v5` + `golang.org/x/crypto/bcrypt`
Python 버전과 동일하게 HS256 JWT(24h 만료, 갱신 없음), bcrypt 해시를 사용한다. `Authorization: Bearer <token>` 헤더 파싱은 chi 미들웨어로 구현.

### 4. API 문서화: `swaggo/swag`
핸들러 함수 위에 swaggo 어노테이션 주석을 작성해 OpenAPI 스펙을 생성하고, `swagger-ui` 정적 자산으로 `/docs`에서 Swagger UI를 제공한다. JWT Bearer 보안 스키마를 OpenAPI에 선언해 Swagger UI의 Authorize 버튼으로 인증 테스트가 가능하도록 한다(Python 버전과 동일한 사용자 경험).

### 5. 에러 응답 표준 이식
Python 버전의 `{ error: { code, message, meta? } }` 형식을 그대로 유지한다. Go에서는 커스텀 `ApiError` 타입(코드/메시지/HTTP 상태/선택적 meta)을 정의하고, 공통 에러 핸들링 미들웨어에서 이를 JSON으로 직렬화한다. 에러 코드 목록(`VALIDATION_ERROR`, `TOO_LONG`, `INVALID_CREDENTIALS`, `TOKEN_EXPIRED`, `FORBIDDEN`, `NOT_OWNER`, `NOT_FOUND`, `EMAIL_TAKEN`, `CONFLICT`)은 Python 버전과 동일하게 유지한다.

### 6. 디렉토리 구조
```
backend/
  cmd/server/main.go        # 엔트리포인트
  internal/
    config/                 # 환경변수, 설정
    db/                     # DB 연결, 마이그레이션
    models/                 # 구조체 정의
    handlers/               # auth.go, teams.go, tasks.go, messages.go
    middleware/             # JWT 인증, CORS, 에러 핸들링
    errors/                 # ApiError 및 표준 에러 생성 헬퍼
  migrations/                # SQL 마이그레이션 파일
  go.mod / go.sum
```

### 7. 테스트: 표준 `testing` + `net/http/httptest`
각 핸들러 그룹(auth/teams/tasks/messages)에 대해 `httptest.NewServer` 또는 `httptest.NewRecorder`를 사용한 통합 테스트를 작성한다. 테스트 DB는 SQLite in-memory(`:memory:`)를 사용해 Python 버전의 pytest 픽스처와 동일한 역할을 하게 한다.

### 8. 마이그레이션 도구
`golang-migrate/migrate` 또는 애플리케이션 기동 시 `CREATE TABLE IF NOT EXISTS` 방식 중 구현 단계에서 확정. MVP 규모(4테이블, 스키마 변경 거의 없음)를 고려하면 후자(기동 시 자동 생성)로 단순화해도 무방.

## Risks / Trade-offs

- [Risk] Go 팀 숙련도가 Python 대비 낮을 경우 초기 개발 속도 저하 → Mitigation: chi/sqlx 등 러닝커브가 낮은 라이브러리 선택, 기존 Python 코드를 참조 구현으로 활용
- [Risk] SQLite와 Postgres 간 SQL 방언 차이(예: `AUTOINCREMENT` vs `SERIAL`, 타임스탬프 타입)로 로컬/운영 동작 불일치 가능 → Mitigation: 표준 SQL 위주로 작성하고 방언별 마이그레이션 파일을 분리, 배포 전 운영 DB 대상 스모크 테스트 필수
- [Risk] Python 버전과 미묘하게 다른 JSON 직렬화(예: datetime 포맷)로 프론트엔드가 깨질 가능성 → Mitigation: ISO 8601 포맷 통일, 기존 Playwright 골든 패스 테스트를 Go 백엔드 대상으로 재실행하여 회귀 확인
- [Risk] swaggo 어노테이션 누락/오류로 Swagger 문서가 불완전해질 수 있음 → Mitigation: 빌드 시 `swag init` 실행을 CI/Makefile에 포함해 어노테이션 오류를 빌드 실패로 드러냄

## Migration Plan

1. Go 프로젝트 스캐폴딩(`go.mod`, 디렉토리 구조) 생성
2. DB 스키마/마이그레이션을 Go 쪽에 동일하게 재정의
3. 엔드포인트를 그룹별(Auth → Team → Kanban → Chat) 순서로 이식하며, 매 그룹마다 Python 버전과 동일한 요청/응답으로 수동 대조
4. 프론트엔드를 Go 백엔드(`localhost:8000`)에 연결해 기존 Playwright 골든 패스를 재실행, 회귀 없음을 확인
5. 전체 이식 완료 후 `backend/`의 Python 코드(venv, requirements.txt, app/, tests/) 삭제
6. 기존 `add-taskflow-mvp`의 보류된 배포 작업(Neon/Vercel)은 Go 바이너리 배포 방식에 맞춰 별도로 재검토

롤백: 이식 도중 문제가 발생하면 Python 백엔드 코드를 삭제하지 않고 그대로 두어(Go 코드는 별도 디렉토리에서 개발 후 마지막에 스왑) 언제든 되돌릴 수 있게 한다.

## Open Questions

- `sqlx` vs `sqlc` 최종 선택은 구현 초반 스파이크 후 확정 (design은 sqlx를 기본값으로 제안)
- DB 마이그레이션을 `golang-migrate`로 관리할지, 앱 기동 시 자동 스키마 생성으로 단순화할지는 구현 단계에서 결정
- Go 백엔드의 Vercel 배포 방식(Go 바이너리를 Vercel Go Runtime으로 배포 vs 별도 호스팅)은 이번 변경 범위 밖이며 이식 완료 후 별도 논의 필요

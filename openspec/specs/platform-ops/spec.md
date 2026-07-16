## ADDED Requirements

### Requirement: 표준 에러 응답 형식
시스템의 모든 4xx/5xx 응답은 `{ error: { code, message, meta? } }` 형태를 따라야 한다(SHALL). `code`는 SCREAMING_SNAKE_CASE 기계 판독 코드여야 하고, `message`는 사용자에게 그대로 노출 가능한 한국어여야 한다.

#### Scenario: 임의의 에러 응답
- **WHEN** 어떤 API 호출이든 4xx 또는 5xx 상태 코드를 반환하면
- **THEN** 응답 바디는 `error.code`와 `error.message` 필드를 포함한다

#### Scenario: 민감정보 비노출
- **WHEN** 로그인 실패 등 인증 관련 에러가 발생하면
- **THEN** 응답 메시지는 이메일 존재 여부 등 민감정보를 드러내지 않는다

### Requirement: API 문서화 (Swagger/OpenAPI)
시스템은 `swaggo/swag` 어노테이션 기반으로 생성된 OpenAPI 문서를 `/docs`(Swagger UI)에서 노출해야 한다(SHALL). JWT Bearer 인증이 필요한 엔드포인트는 OpenAPI 보안 스키마로 선언되어 Swagger UI에서 토큰 입력 후 바로 호출 테스트가 가능해야 한다.

#### Scenario: Swagger UI에서 인증 API 호출
- **WHEN** 개발자가 `/docs`에서 "Authorize" 버튼으로 JWT를 입력하고 보호된 엔드포인트를 실행하면
- **THEN** Swagger UI는 `Authorization: Bearer <token>` 헤더를 자동으로 포함하여 요청을 보내고 실제 API 응답을 표시한다

#### Scenario: 빌드 시 OpenAPI 스펙 생성
- **WHEN** 개발자가 빌드 파이프라인(`swag init` 포함)을 실행하면
- **THEN** 핸들러 어노테이션으로부터 최신 OpenAPI 스펙 파일이 생성되고 `/docs`에 반영된다

### Requirement: 환경 분리 (로컬/운영)
시스템은 `DATABASE_URL` 환경변수 하나로 로컬 SQLite와 운영 Neon Postgres 간 전환이 가능해야 한다(SHALL). 로컬 SQLite 파일은 버전관리에서 제외되어야 한다.

#### Scenario: 로컬 실행
- **WHEN** `DATABASE_URL`이 설정되지 않거나 SQLite 경로를 가리키면
- **THEN** 애플리케이션은 로컬 SQLite 파일을 사용해 기동한다

#### Scenario: 운영 배포
- **WHEN** `DATABASE_URL`이 Neon Postgres 접속 문자열로 설정되면
- **THEN** 애플리케이션은 동일한 코드로 Postgres에 연결하여 기동한다

### Requirement: 핵심 플로우 자동 테스트
시스템은 회원가입, 로그인, 칸반 태스크 CRUD(생성/상태변경/삭제 권한 포함), 채팅 메시지 송수신(폴링 포함)에 대한 자동화된 테스트를 Go 표준 `testing` 패키지와 `net/http/httptest`로 갖춰야 한다(SHALL). 전면적인 테스트 자동화는 요구하지 않으며, 그 외 시나리오는 수동 동작 확인으로 검증한다.

#### Scenario: 핵심 플로우 테스트 실행
- **WHEN** 개발자가 `go test ./...`를 실행하면
- **THEN** 회원가입/로그인/칸반 CRUD/채팅 송수신에 대한 테스트가 실행되고 통과한다

### Requirement: 배포 파이프라인
시스템은 GitHub의 `main` 브랜치 push 시 Vercel을 통해 프론트엔드(정적 파일)와 백엔드(Serverless Functions)가 함께 자동 배포되어야 한다(SHALL).

#### Scenario: main 브랜치 배포
- **WHEN** 코드가 GitHub `main` 브랜치에 push되면
- **THEN** Vercel은 프론트엔드와 백엔드를 함께 빌드하고 배포한다

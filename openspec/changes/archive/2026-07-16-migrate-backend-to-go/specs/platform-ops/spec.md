## MODIFIED Requirements

### Requirement: API 문서화 (Swagger/OpenAPI)
시스템은 `swaggo/swag` 어노테이션 기반으로 생성된 OpenAPI 문서를 `/docs`(Swagger UI)에서 노출해야 한다(SHALL). JWT Bearer 인증이 필요한 엔드포인트는 OpenAPI 보안 스키마로 선언되어 Swagger UI에서 토큰 입력 후 바로 호출 테스트가 가능해야 한다.

#### Scenario: Swagger UI에서 인증 API 호출
- **WHEN** 개발자가 `/docs`에서 "Authorize" 버튼으로 JWT를 입력하고 보호된 엔드포인트를 실행하면
- **THEN** Swagger UI는 `Authorization: Bearer <token>` 헤더를 자동으로 포함하여 요청을 보내고 실제 API 응답을 표시한다

#### Scenario: 빌드 시 OpenAPI 스펙 생성
- **WHEN** 개발자가 빌드 파이프라인(`swag init` 포함)을 실행하면
- **THEN** 핸들러 어노테이션으로부터 최신 OpenAPI 스펙 파일이 생성되고 `/docs`에 반영된다

### Requirement: 핵심 플로우 자동 테스트
시스템은 회원가입, 로그인, 칸반 태스크 CRUD(생성/상태변경/삭제 권한 포함), 채팅 메시지 송수신(폴링 포함)에 대한 자동화된 테스트를 Go 표준 `testing` 패키지와 `net/http/httptest`로 갖춰야 한다(SHALL). 전면적인 테스트 자동화는 요구하지 않으며, 그 외 시나리오는 수동 동작 확인으로 검증한다.

#### Scenario: 핵심 플로우 테스트 실행
- **WHEN** 개발자가 `go test ./...`를 실행하면
- **THEN** 회원가입/로그인/칸반 CRUD/채팅 송수신에 대한 테스트가 실행되고 통과한다

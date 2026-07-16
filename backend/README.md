# TaskFlow Backend (Go)

FastAPI(Python)로 구현되었던 TaskFlow MVP 백엔드를 Go로 재구현한 버전입니다. API 계약(엔드포인트, 요청/응답, 에러 코드)은 이전 Python 버전과 동일합니다.

## 로컬 실행

```bash
go run ./cmd/server
```

기본적으로 `./taskflow.db` (SQLite)를 사용하며 `:8000` 포트에서 기동합니다.

## 테스트

```bash
go test ./...
```

## API 문서 (Swagger)

서버 실행 후 브라우저에서 `http://localhost:8000/docs` 접속. JWT 발급 후 "Authorize" 버튼으로 토큰을 입력하면 인증이 필요한 API도 바로 호출 테스트 가능합니다.

Swagger UI 정적 자산(`swagger-ui.css`, `swagger-ui-bundle.js`)은 `cmd/server/swaggerui/`에 벤더링되어 `go:embed`로 바이너리에 내장됩니다 — unpkg 등 외부 CDN 없이 완전히 오프라인에서 동작합니다.

OpenAPI 스펙을 어노테이션으로부터 재생성하려면:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

## 환경 변수

| 변수 | 설명 | 기본값 |
|---|---|---|
| `DATABASE_URL` | `sqlite:///./taskflow.db` 또는 `postgres://...` (docker-compose의 로컬 `postgres` 컨테이너, 또는 원한다면 Neon 등 외부 Postgres) | `sqlite:///./taskflow.db` |
| `JWT_SECRET` | JWT 서명 시크릿 | `dev-secret-change-in-production` (운영에서는 반드시 변경) |
| `CORS_ORIGINS` | 콤마로 구분된 허용 도메인 목록 | `http://localhost:8000,http://localhost:5500` |
| `PORT` | 서버 포트 | `8000` |

## 디렉토리 구조

```
backend/
  cmd/server/       # 엔트리포인트, Swagger UI 라우트
    swaggerui/      # 벤더링된 swagger-ui-dist 정적 자산 (go:embed, CDN 미사용)
  internal/
    app/            # 공통 App(DB, Settings) 구조체
    config/         # 환경변수 로딩
    db/             # DB 연결, 스키마
    models/         # 구조체 정의
    handlers/       # auth, teams, tasks, messages 핸들러
    middleware/      # JWT 인증, CORS, 에러 응답, 보안 유틸
    errors/         # 표준 에러 타입/헬퍼
    server/         # 라우터 조립 (main.go와 테스트 공용)
    testsupport/    # 테스트용 in-memory DB 헬퍼
  docs/             # swag init으로 생성된 OpenAPI 스펙
```

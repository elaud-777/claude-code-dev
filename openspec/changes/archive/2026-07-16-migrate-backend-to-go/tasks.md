## 1. Go 프로젝트 스캐폴딩

- [x] 1.1 `go mod init` 실행, 디렉토리 구조 생성 (`cmd/server`, `internal/config`, `internal/db`, `internal/models`, `internal/handlers`, `internal/middleware`, `internal/errors`)
- [x] 1.2 의존성 추가: chi, sqlx, golang-jwt/jwt/v5, golang.org/x/crypto/bcrypt, sqlite3/postgres 드라이버, swaggo/swag
- [x] 1.3 `internal/config`: `DATABASE_URL`, `JWT_SECRET`, `CORS_ORIGINS` 환경변수 로딩 구현
- [x] 1.4 `internal/db`: SQLite/Postgres 겸용 DB 연결 초기화 (스킴 기반 드라이버 분기)

## 2. DB 스키마

- [x] 2.1 `users`, `teams`, `tasks`, `messages` 테이블 생성 SQL 작성 (기동 시 `CREATE TABLE IF NOT EXISTS` 또는 마이그레이션 도구, design.md 결정 사항 반영)
- [x] 2.2 인덱스 추가: tasks(team_id, created_at), messages(team_id, created_at), teams(invite_code) unique, users(team_id)
- [x] 2.3 Go 구조체(models) 정의: User, Team, Task, Message (Python SQLAlchemy 모델과 필드 1:1 대응)

## 3. 공통 인프라

- [x] 3.1 `internal/errors`: ApiError 타입 및 표준 에러 생성 헬퍼 (VALIDATION_ERROR, TOO_LONG, INVALID_CREDENTIALS, TOKEN_EXPIRED, FORBIDDEN, NOT_OWNER, NOT_FOUND, EMAIL_TAKEN, CONFLICT)
- [x] 3.2 에러 핸들링 미들웨어: ApiError를 `{ error: { code, message, meta? } }` JSON으로 직렬화
- [x] 3.3 CORS 미들웨어 설정 (허용 도메인 환경변수 기반)
- [x] 3.4 JWT 인증 미들웨어: `Authorization: Bearer` 파싱, 만료/검증 실패 시 401 TOKEN_EXPIRED
- [x] 3.5 팀 멤버십 검증 헬퍼 (모든 `/teams/{id}/*` 라우트에서 재사용)
- [x] 3.6 swaggo 어노테이션 기반 `/docs` Swagger UI 및 JWT Bearer 보안 스키마 설정

## 4. 인증 핸들러 (auth)

- [x] 4.1 `POST /auth/signup`: 이메일 형식 검증, 중복 확인(409 EMAIL_TAKEN), 비밀번호 8자 이상 검증, bcrypt 해시, JWT 발급
- [x] 4.2 `POST /auth/login`: 자격 증명 검증, 실패 시 401 INVALID_CREDENTIALS, 성공 시 24h JWT 발급
- [x] 4.3 `POST /auth/logout`: stateless, 인증된 요청에 200만 반환
- [x] 4.4 `GET /auth/me`: JWT 검증 후 현재 사용자 정보 반환
- [x] 4.5 Python 버전과 응답 바디 대조 테스트 (동일 요청 → 동일 JSON 구조)

## 5. 팀 관리 핸들러 (team-management)

- [x] 5.1 초대코드 생성기: `^[A-Z]{4}-[0-9]{4}$` 형식, 충돌 시 재시도
- [x] 5.2 `POST /teams`: 팀 생성, owner_id 지정, team_id 갱신
- [x] 5.3 `POST /teams/join`: 형식(400)/미존재(404)/이미 소속(409) 에러 분기
- [x] 5.4 `GET /teams/{id}`: 멤버십 검증 후 팀 정보 반환
- [x] 5.5 `GET /teams/{id}/members`: 멤버 목록(owner 표시 포함) 반환
- [x] 5.6 `DELETE /teams/{id}/leave`: 본인 team_id를 NULL로 갱신

## 6. 칸반 핸들러 (kanban)

- [x] 6.1 `POST /teams/{id}/tasks`: 제목(1-100자) 검증, status=TODO 기본값, assignee_id nullable
- [x] 6.2 `GET /teams/{id}/tasks`: 전체/@me/미할당 필터, 생성일 정렬
- [x] 6.3 `GET /tasks/{id}`: 단건 상세 조회
- [x] 6.4 `PUT /tasks/{id}`: 제목/담당자 수정
- [x] 6.5 `PATCH /tasks/{id}/status`: 상태 변경 전용
- [x] 6.6 `DELETE /tasks/{id}`: creator 또는 team owner만 허용, 그 외 403 FORBIDDEN

## 7. 채팅 핸들러 (team-chat)

- [x] 7.1 `POST /teams/{id}/messages`: 1000자 제한 서버측 검증(400 TOO_LONG)
- [x] 7.2 `GET /teams/{id}/messages`: `since=` 파라미터 지원, 없으면 최근 50개
- [x] 7.3 `DELETE /messages/{id}`: 작성자 본인만 허용(owner 포함 타인은 403 NOT_OWNER)

## 8. 테스트

- [x] 8.1 테스트 헬퍼: SQLite in-memory 기반 httptest 서버 셋업
- [x] 8.2 회원가입 테스트: 성공/이메일 중복/약한 비밀번호
- [x] 8.3 로그인 테스트: 성공/실패
- [x] 8.4 칸반 CRUD 테스트: 생성/상태변경/삭제 권한(creator·owner·타인)
- [x] 8.5 채팅 테스트: 메시지 전송/폴링(since=)/삭제 권한
- [x] 8.6 `go test ./...` 전체 통과 확인

## 9. 프론트엔드 연동 회귀 검증

- [x] 9.1 프론트엔드(`frontend/`)를 Go 백엔드(`localhost:8000`)에 연결해 기동
- [x] 9.2 기존 Playwright 골든 패스 재실행: 회원가입 → 팀 생성 → 태스크 생성/상태변경 → 채팅 전송 → 멤버 목록 조회
- [x] 9.3 모바일 반응형(햄버거 메뉴, 칸반 컬럼) 회귀 확인

## 10. 정리

- [x] 10.1 Python 백엔드 코드 삭제 (`backend/app`, `backend/tests`, `backend/requirements.txt`, `backend/.venv`)
- [x] 10.2 `backend/.gitignore`를 Go 프로젝트 기준으로 갱신 (`*.db`, 빌드 산출물 등)
- [x] 10.3 README/실행 안내를 Go 기준으로 갱신 (`go run ./cmd/server` 등)

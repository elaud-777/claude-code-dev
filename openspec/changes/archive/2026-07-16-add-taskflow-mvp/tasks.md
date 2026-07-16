## 1. 백엔드 프로젝트 셋업

- [x] 1.1 `backend/` 디렉토리 생성, FastAPI + SQLAlchemy + pydantic + python-jose + bcrypt(passlib) 의존성 설정 (requirements.txt / pyproject.toml)
- [x] 1.2 로컬 SQLite(`sqlite:///./taskflow.db`) / 운영 Neon Postgres를 `DATABASE_URL` 환경변수로 전환하는 DB 연결 모듈 작성
- [x] 1.3 CORS 허용 도메인 설정 및 에러 응답 표준 `{ error: { code, message, meta? } }` 공통 핸들러 작성
- [x] 1.4 FastAPI OpenAPI 보안 스키마(HTTPBearer) 설정하여 `/docs` Swagger UI에서 JWT 인증 테스트 가능하도록 구성

## 2. DB 스키마

- [x] 2.1 `users` 테이블 정의: id, email(unique), password_hash, team_id(FK, nullable), created_at
- [x] 2.2 `teams` 테이블 정의: id, name, invite_code(unique), owner_id(FK), created_at
- [x] 2.3 `tasks` 테이블 정의: id, team_id(FK), title, status(TODO/DOING/DONE), creator_id(FK), assignee_id(FK, nullable), created_at
- [x] 2.4 `messages` 테이블 정의: id, team_id(FK), user_id(FK), content, created_at
- [x] 2.5 인덱스 추가: tasks(team_id, created_at), messages(team_id, created_at), teams(invite_code) unique, users(team_id)

## 3. 인증 (auth capability)

- [x] 3.1 `POST /auth/signup` 구현: 이메일 형식 검증, 중복 확인(409 EMAIL_TAKEN), 비밀번호 8자 이상 검증, bcrypt 해시 저장, JWT 발급
- [x] 3.2 `POST /auth/login` 구현: 자격 증명 검증, 실패 시 통일된 401 INVALID_CREDENTIALS, 성공 시 24h JWT 발급
- [x] 3.3 `POST /auth/logout` 구현: stateless, 인증된 요청에 200만 반환
- [x] 3.4 `GET /auth/me` 구현: JWT 검증 후 현재 사용자 정보 반환
- [x] 3.5 JWT 검증 미들웨어/의존성 작성: 만료 시 401 TOKEN_EXPIRED

## 4. 팀 관리 (team-management capability)

- [x] 4.1 초대코드 생성기 작성: `^[A-Z]{4}-[0-9]{4}$` 형식, 충돌 시 재시도
- [x] 4.2 `POST /teams` 구현: 팀 생성, owner_id 지정, 생성자 team_id 갱신
- [x] 4.3 `POST /teams/join` 구현: 형식(400)/미존재(404)/이미 소속(409) 에러 분기 및 정상 합류 처리
- [x] 4.4 `GET /teams/{id}` 구현: 멤버십 검증(403 FORBIDDEN) 후 팀 정보 반환
- [x] 4.5 `GET /teams/{id}/members` 구현: 멤버십 검증 후 멤버 목록(owner 표시 포함) 반환
- [x] 4.6 `DELETE /teams/{id}/leave` 구현: 본인 team_id를 NULL로 갱신
- [x] 4.7 팀 멤버십 검증 공통 의존성 작성 (모든 `/teams/{id}/*` 라우트에 재사용)

## 5. 칸반 (kanban capability)

- [x] 5.1 `POST /teams/{id}/tasks` 구현: 제목(1-100자) 검증, status=TODO 기본값, assignee_id nullable 처리
- [x] 5.2 `GET /teams/{id}/tasks` 구현: 전체/`@me`(assignee_id 기준)/미할당 필터, 생성일 정렬
- [x] 5.3 `GET /tasks/{id}` 구현: 단건 상세 조회
- [x] 5.4 `PUT /tasks/{id}` 구현: 제목/담당자 수정 (상태 변경 제외)
- [x] 5.5 `PATCH /tasks/{id}/status` 구현: 드래그로 인한 상태 변경 전용
- [x] 5.6 `DELETE /tasks/{id}` 구현: creator 또는 team owner만 허용, 그 외 403 FORBIDDEN

## 6. 채팅 (team-chat capability)

- [x] 6.1 `POST /teams/{id}/messages` 구현: 1000자 제한 서버측 검증(400 TOO_LONG)
- [x] 6.2 `GET /teams/{id}/messages` 구현: `since=` 파라미터 지원, 없으면 최근 50개 반환
- [x] 6.3 `DELETE /messages/{id}` 구현: 작성자 본인만 허용(owner 포함 타인은 403 NOT_OWNER)

## 7. 백엔드 테스트

- [x] 7.1 pytest 셋업 (테스트용 SQLite in-memory 또는 임시 파일 DB)
- [x] 7.2 회원가입 테스트: 성공/이메일 중복/약한 비밀번호
- [x] 7.3 로그인 테스트: 성공/실패
- [x] 7.4 칸반 CRUD 테스트: 생성/상태변경/삭제 권한(creator·owner·타인)
- [x] 7.5 채팅 테스트: 메시지 전송/폴링(since=)/삭제 권한

## 8. 프론트엔드 프로젝트 셋업

- [x] 8.1 `frontend/` 디렉토리 생성, Vanilla JS + Tailwind(CDN 또는 빌드) 셋업
- [x] 8.2 API 클라이언트 모듈 작성: fetch 래퍼, JWT 헤더 자동 첨부, 401 시 로그인 화면 redirect
- [x] 8.3 JWT localStorage 저장/조회/삭제 유틸 작성

## 9. 프론트엔드 화면 구현

- [x] 9.1 로그인/회원가입 화면 (초기/입력중/처리중/에러 상태 포함)
- [x] 9.2 팀 선택 화면 (팀 만들기 + 초대코드 합류, team_id NULL 강제 라우팅)
- [x] 9.3 칸반 화면: 3컬럼, 카드 생성(인라인 입력), 드래그 상태 변경, 상세/수정 모달, 삭제 확인 다이얼로그, 담당자 필터
- [x] 9.4 채팅 화면: 메시지 목록, 5초 폴링(since=), 1000자 카운터, 본인 메시지 삭제
- [x] 9.5 팀 멤버 목록 사이드 패널
- [x] 9.6 반응형 대응: 모바일 브레이크포인트(칸반 스와이프, 채팅 풀스크린, 햄버거 메뉴)
- [x] 9.7 공통 에러 처리: 403/404/409/400 각 케이스 사용자 메시지 표시, JWT 만료 시 자동 로그인 화면 이동

## 10. 배포 (DEFERRED — 사용자 Neon/Vercel 계정 연결 필요로 보류, 별도 작업으로 진행 예정)

- [ ] 10.1 ~~Neon Postgres 프로젝트 생성 및 `DATABASE_URL` 확보~~ (deferred)
- [ ] 10.2 ~~Vercel 프로젝트 연결: frontend(정적 파일) + backend(Serverless Functions)~~ (deferred)
- [ ] 10.3 ~~Vercel 환경변수에 `DATABASE_URL`, CORS 허용 도메인 등록~~ (deferred)
- [ ] 10.4 ~~배포 후 스모크 테스트: 회원가입 → 로그인 → 팀 생성 → 칸반 카드 생성/이동 → 채팅 전송 수동 확인~~ (deferred)

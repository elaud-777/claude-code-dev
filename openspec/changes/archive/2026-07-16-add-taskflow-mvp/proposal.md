## Why

소규모 팀(3-5인)이 칸반 보드와 팀 채팅을 하나의 화면에서 오갈 수 있는 도구가 없어 진행 상황 파악과 합의에 시간이 걸린다. TaskFlow MVP는 팀 리더/팀원/신규 합류자가 한 화면에서 태스크 진행과 짧은 채팅 합의를 동시에 처리하도록 하여, 신규 합류자가 1분 안에 컨텍스트를 파악하고 팀이 어디서나(PC/모바일) 접근할 수 있게 한다. Day 2 마무리 시점까지 완성 가능한 범위로 기능을 고정해 지금 시작한다.

## What Changes

- 이메일/비밀번호 기반 회원가입·로그인 추가, JWT(24h 만료, 갱신 없음) 발급 및 bcrypt 비밀번호 해시 적용
- 팀 생성 및 초대코드(`AAAA-9999` 형식) 발급/합류/멤버 목록 조회 추가. 사용자는 동시에 하나의 팀에만 소속(`users.team_id`, nullable)
- 칸반 보드 추가: TODO/DOING/DONE 3컬럼, 카드 생성/조회/제목수정/상태변경(드래그)/삭제, 담당자(`assignee_id`, nullable) 지정 및 필터(전체/@me/미할당)
- 팀 단위 채팅 추가: 5초 폴링 기반 메시지 송수신(`since=` 파라미터로 증분 조회), 1000자 제한, 본인 메시지만 삭제 가능
- FastAPI 기반 REST API 18종 구현 및 FastAPI 기본 제공 Swagger(OpenAPI `/docs`)를 JWT Bearer 인증 테스트 가능하도록 노출
- 핵심 플로우(회원가입/로그인/칸반 CRUD/채팅 송수신)에 대한 기본 자동 테스트(pytest) 작성 — 전면 자동화는 아니며 수동 동작 확인이 우선
- 로컬(SQLite)/운영(Vercel + Neon Postgres) 배포 환경 분리, `DATABASE_URL` 환경변수로 전환
- 에러 응답 표준 `{ error: { code, message } }` 전 API 공통 적용
- 표준 대비 변경: `PATCH /tasks/{id}/status`로 상태변경을 `PUT /tasks/{id}`(제목수정)와 분리, `GET /messages/{id}` 대신 `GET /teams/{id}`·`DELETE /teams/{id}/leave` 추가

## Capabilities

### New Capabilities
- `auth`: 회원가입, 로그인, 로그아웃(stateless), 현재 사용자 조회, JWT 발급/검증, bcrypt 해시
- `team-management`: 팀 생성, 초대코드 발급/합류, 멤버 목록 조회, 팀 나가기, 멤버십 기반 접근 제어(403)
- `kanban`: 태스크 생성/조회/제목수정/상태변경/삭제, 담당자 지정 및 필터/정렬, 권한 기반 삭제(creator 또는 owner)
- `team-chat`: 팀 단위 메시지 송수신(폴링), 1000자 제한, 본인 메시지 삭제
- `platform-ops`: 공통 에러 응답 표준, API 문서화(Swagger/OpenAPI), 로컬/운영 환경 분리 및 배포(Vercel + Neon), 핵심 플로우 자동 테스트

### Modified Capabilities
(없음 — 모두 신규 프로젝트이므로 기존 스펙 없음)

## Impact

- 신규 코드베이스: `D:\claude-code-dev\backend`(FastAPI + SQLAlchemy), `D:\claude-code-dev\frontend`(Vanilla JS + Tailwind)
- 신규 DB 스키마 4테이블: `users`, `teams`, `tasks`, `messages` (로컬 SQLite / 운영 Neon Postgres)
- 신규 API 18개 엔드포인트 (Auth 4 + Team 5 + Task 6 + Chat 3)
- 신규 배포 파이프라인: GitHub → Vercel(FE+BE) + Neon(DB), `DATABASE_URL` 환경변수로 로컬/운영 전환
- 외부 의존성 없음(기존 시스템과 통합 없는 신규 프로젝트)

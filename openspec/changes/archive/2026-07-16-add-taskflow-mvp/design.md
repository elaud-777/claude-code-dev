## Context

TaskFlow MVP는 프로그램정의서(미션/페르소나/기능 5종/1차 API·DB 초안)와 스토리보드(42슬라이드, 화면 상태·에러 케이스·결정추적표 8건 반영한 통합본)를 입력으로 한다. 두 문서가 일부(API 구성 등)에서 다를 경우 스토리보드를 진본으로 채택한다. 신규 프로젝트이며 기존 시스템과의 통합은 없다. 개발/검증은 로컬(FastAPI + SQLite)에서 진행하고 운영은 Vercel(FE+BE) + Neon(Postgres)로 배포한다.

제약:
- 팀 규모 3-5인, 동시 접속 5명/팀 이하 가정
- 한국어 UI만, 단일 시간대(KST) 가정
- Day 2 마무리 시점까지 완성 가능한 범위로 기능 고정

## Goals / Non-Goals

**Goals:**
- 인증 → 팀 합류 → 칸반 작업 → 채팅 합의로 이어지는 핵심 사용자 여정을 끊김 없이 구현
- 스토리보드 결정추적표 8건을 그대로 반영해 API/DB 스키마의 모호함을 제거
- 로컬 개발과 운영 배포 간 전환을 `DATABASE_URL` 하나로 단순화
- Swagger(OpenAPI) 문서에서 JWT 인증 API를 바로 테스트 가능하게 함
- 핵심 플로우(회원가입/로그인/칸반 CRUD/채팅)에 대한 최소 자동 테스트 확보

**Non-Goals:**
- 알림(이메일/SMS/푸시), 파일 첨부, 전문 검색, 세분화된 권한(페이지별), 다국어, WebSocket 실시간 통신 — 모두 범위 외
- 전면적인 테스트 자동화(커버리지 목표 없음) — 수동 확인이 기본, 핵심 플로우만 자동 테스트
- 알림 갱신 토큰, 팀 추방/역할 변경, 초대코드 재발급 — Day 2 범위 외

## Decisions

### 1. `users.team_id` (nullable)로 멤버십 표현 (결정 #1)
`owner_id`만으로는 일반 멤버가 어느 팀에 속하는지 표현할 수 없다. "1인 1팀" 가정 하에 `users.team_id`를 추가한다. `team_id IS NULL`이면 팀 미가입 상태로 취급하고, 모든 `/teams/*` 라우트 진입 전에 이 값을 검사해 강제 리다이렉트(팀 선택 화면)한다.
- 대안 검토: 별도 `team_members` N:M 테이블 — "1인 1팀" 가정에서는 불필요한 복잡도이므로 기각.

### 2. 상태 변경과 제목 수정 API 분리 (결정 #3)
`PUT /tasks/{id}`(제목 수정)과 `PATCH /tasks/{id}/status`(상태 변경, 드래그 앤 드롭)를 분리한다. REST 의미상 PATCH는 부분 갱신에 적합하고, 프론트엔드에서 드래그 이벤트와 상세 모달 저장 이벤트가 서로 다른 API를 호출하도록 명확히 분리할 수 있다.

### 3. `tasks.assignee_id` (nullable) 추가 (결정 #4)
`creator_id`(생성자) ≠ "내 태스크" 이므로 별도 `assignee_id`를 추가한다. "내 태스크" 필터는 `WHERE assignee_id = current_user_id`로 정의하며 `creator_id`는 사용하지 않는다. `assignee_id IS NULL`인 카드는 "미할당"으로 표시한다.

### 4. Stateless 로그아웃 (결정 #5)
JWT는 stateless이므로 서버에 블랙리스트를 두지 않는다. `POST /auth/logout`은 인증된 요청이면 `200 {}`만 반환하고, 실제 토큰 폐기는 클라이언트의 `localStorage.removeItem`으로 처리한다. 트레이드오프: 토큰이 탈취되어도 만료(24h) 전까지 강제 무효화할 수 없음 — Day 2 범위에서는 감수.

### 5. 권한 모델: 멤버십 403 + 소유권 기반 삭제 (결정 #6)
- 모든 `/teams/{id}/*` 라우트는 JWT 검증 후 `user.team_id == {id}`를 확인. 불일치 시 403 `FORBIDDEN`.
- 태스크 삭제: `creator_id == current_user` 또는 `team.owner_id == current_user`만 허용(owner는 오버라이드 가능).
- 메시지 삭제: `user_id == current_user`만 허용 — owner도 타인 메시지는 삭제 불가(커뮤니티 신뢰 모델).

### 6. 성능 지표는 정성 검증 (결정 #7)
드래그 반응 50ms, API 응답 100ms, 신규 합류자 1분 파악 등은 자동 측정 도구를 두지 않고 수동 확인으로 검증한다. 자동 성능 테스트 도구 도입은 범위 외.

### 7. API 구성 최종본 채택 (결정 #8)
`GET /messages/{id}`(용도 모호)를 제거하고 `GET /teams/{id}`(팀 정보 조회), `DELETE /teams/{id}/leave`(팀 나가기)로 대체한다. 최종 API 구성은 Auth 4 + Team 5 + Task 6 + Chat 3 = 18개(스토리보드 G·02 기준).

### 8. 백엔드: FastAPI + SQLAlchemy, ORM으로 SQLite/Postgres 양쪽 호환
로컬은 `sqlite:///./taskflow.db`, 운영은 Neon의 `DATABASE_URL`(postgres://...)을 그대로 주입한다. SQLAlchemy를 사용해 두 DB 방언 차이를 추상화한다. 라이브러리 선택(SQLAlchemy vs Tortoise 등), 디렉토리 구조, 마이그레이션 도구는 구현 단계에서 결정.

### 9. 프론트엔드: Vanilla JS + Tailwind, 프레임워크 없음
사용자가 확정한 스택. SPA/MPA 여부, 라우팅 방식, CSS 빌드 방식(CDN vs 빌드), 상태 관리 방식은 구현 단계에서 결정.

### 10. Swagger/OpenAPI 노출
FastAPI가 기본 제공하는 `/docs`(Swagger UI), `/redoc`을 그대로 사용한다. JWT Bearer 인증이 필요한 엔드포인트는 OpenAPI 보안 스키마(`HTTPBearer`)로 선언해 Swagger UI의 "Authorize" 버튼으로 토큰을 넣고 바로 호출 테스트가 가능하게 한다. 별도 Swagger 설치는 불필요.

### 11. 테스트 범위: 핵심 플로우만 자동화
pytest로 다음을 커버: 회원가입(성공/이메일 중복/약한 비밀번호), 로그인(성공/실패), 칸반 카드 생성·상태변경·삭제(권한 포함), 채팅 메시지 송수신(폴링 `since=` 포함). 프론트엔드 자동 테스트는 범위 외(수동 확인). 전체 커버리지 목표는 없음 — 이 핵심 플로우 외 확장은 이번 변경 범위에 포함하지 않는다.

### 12. 에러 응답 표준
모든 4xx/5xx 응답은 `{ error: { code: SCREAMING_SNAKE, message: 한국어, meta?: {...} } }` 형태로 통일한다. 주요 코드: `VALIDATION_ERROR`(400), `TOO_LONG`(400), `INVALID_CREDENTIALS`(401), `TOKEN_EXPIRED`(401), `FORBIDDEN`(403), `NOT_OWNER`(403), `NOT_FOUND`(404), `EMAIL_TAKEN`(409).

## Risks / Trade-offs

- [Risk] JWT 갱신 토큰 없음 → 24h마다 강제 재로그인, 장시간 세션 사용자 불편 → Mitigation: Day 2 범위 내 의도된 트레이드오프로 명시, 추후 개선 과제로 이월
- [Risk] 5초 폴링 방식은 팀원이 늘어나면 서버 부하 증가 → Mitigation: 팀 규모 5명 이하 가정으로 범위 고정, WebSocket 전환은 향후 과제
- [Risk] `assignee_id`/`team_id` 등 nullable 컬럼이 늘어나며 쿼리 조건 분기 증가 → Mitigation: 필터 로직을 한 곳(서비스 레이어)에 모아 산발적 분기 방지
- [Risk] SQLite(로컬)와 Postgres(운영) 간 동작 차이(예: 타입 강제, 동시성) → Mitigation: SQLAlchemy Core/ORM 표준 타입만 사용, 로컬에서 발견 못한 이슈는 운영 배포 후 스모크 테스트로 확인
- [Risk] 초대코드 형식(`^[A-Z]{4}-[0-9]{4}$`) 충돌(중복 생성) → Mitigation: DB UNIQUE 제약 + 생성 시 재시도 로직

## Migration Plan

신규 프로젝트이므로 기존 데이터 마이그레이션은 없다. 배포 순서:
1. 로컬에서 FastAPI + SQLite로 전체 기능 구현 및 수동/자동 테스트 통과 확인
2. GitHub 저장소에 push
3. Vercel에 backend(Serverless Functions)와 frontend(정적 파일) 연결
4. Neon Postgres 프로비저닝 후 `DATABASE_URL` 환경변수를 Vercel에 등록
5. 배포 후 스모크 테스트: 회원가입 → 로그인 → 팀 생성 → 칸반 카드 생성/이동 → 채팅 전송까지 1회 수동 확인

롤백: Vercel의 이전 배포로 즉시 롤백 가능(Vercel 기본 기능). DB 마이그레이션이 없으므로 스키마 롤백 이슈 없음.

## Open Questions

- Neon Postgres의 pooled connection 문자열을 Vercel Serverless Functions에서 사용할 때 connection 재사용 전략(SQLAlchemy pool 설정)은 구현 단계에서 확정 필요
- 초대코드 자동 생성 시 팀 이름 기반 접두어를 쓸지, 완전 랜덤 문자열을 쓸지는 구현 단계에서 결정(스토리보드 예시는 팀 이름 기반 "FRNT-2026")

## ADDED Requirements

### Requirement: 회원가입
시스템은 이메일과 비밀번호로 신규 계정을 생성해야 한다(SHALL). 이메일은 형식 검증(RFC 형식)을 거치고 중복 여부를 확인해야 하며, 비밀번호는 8자 이상이어야 하고 bcrypt로 해시되어 저장되어야 한다(MUST). 가입 성공 시 즉시 JWT를 발급하여 이메일 인증 절차 없이 계정을 활성화해야 한다(SHALL).

#### Scenario: 정상 가입
- **WHEN** 사용자가 유효한 이메일과 8자 이상 비밀번호로 `POST /auth/signup`을 호출하면
- **THEN** 시스템은 `users` 테이블에 레코드를 생성하고 `201 Created`와 JWT를 반환한다

#### Scenario: 이메일 형식 오류
- **WHEN** 사용자가 `user@invalid`처럼 이메일 형식이 아닌 값으로 가입을 시도하면
- **THEN** 시스템은 `400 { error: { code: 'VALIDATION_ERROR' } }`를 반환한다

#### Scenario: 이메일 중복
- **WHEN** 사용자가 이미 가입된 이메일로 가입을 시도하면
- **THEN** 시스템은 `409 { error: { code: 'EMAIL_TAKEN' } }`를 반환한다

#### Scenario: 비밀번호 약함
- **WHEN** 사용자가 8자 미만 비밀번호로 가입을 시도하면
- **THEN** 시스템은 `400 { error: { code: 'VALIDATION_ERROR' } }`를 반환한다

### Requirement: 로그인
시스템은 이메일과 비밀번호를 검증하고 성공 시 24시간 유효한 JWT를 발급해야 한다(SHALL). 실패 시 이메일 존재 여부를 노출하지 않는 단일 에러 메시지를 반환해야 한다(MUST).

#### Scenario: 정상 로그인
- **WHEN** 사용자가 올바른 이메일/비밀번호로 `POST /auth/login`을 호출하면
- **THEN** 시스템은 `200 OK`와 함께 JWT(만료 24h), `user.id`, `user.email`, `user.team_id`를 반환한다

#### Scenario: 자격 증명 오류
- **WHEN** 사용자가 존재하지 않는 이메일 또는 틀린 비밀번호로 로그인을 시도하면
- **THEN** 시스템은 이메일 존재 여부와 무관하게 동일한 `401 { error: { code: 'INVALID_CREDENTIALS' } }`를 반환한다

### Requirement: 로그아웃 (stateless)
시스템은 JWT를 stateless로 취급하여 블랙리스트를 유지하지 않아야 한다(SHALL NOT). `POST /auth/logout`은 인증된 요청에 대해 `200`만 반환하며, 실제 토큰 폐기는 클라이언트 책임이다.

#### Scenario: 로그아웃 호출
- **WHEN** 인증된 사용자가 `POST /auth/logout`을 호출하면
- **THEN** 시스템은 `200 {}`를 반환하고 서버 측에서 토큰을 무효화하지 않는다

### Requirement: 현재 사용자 조회
시스템은 유효한 JWT로 요청한 사용자에게 자신의 계정 정보(id, email, team_id)를 반환해야 한다(SHALL).

#### Scenario: 인증된 사용자 조회
- **WHEN** 유효한 JWT를 가진 사용자가 `GET /auth/me`를 호출하면
- **THEN** 시스템은 `200`과 함께 해당 사용자의 id, email, team_id를 반환한다

### Requirement: JWT 만료 처리
시스템은 만료되었거나 유효하지 않은 JWT로 보호된 리소스에 접근 시 `401 TOKEN_EXPIRED`를 반환해야 한다(SHALL). 갱신 토큰은 제공하지 않는다(Non-Goal).

#### Scenario: 만료된 토큰으로 API 호출
- **WHEN** 발급 후 24시간이 지난 JWT로 보호된 엔드포인트를 호출하면
- **THEN** 시스템은 `401 { error: { code: 'TOKEN_EXPIRED' } }`를 반환한다

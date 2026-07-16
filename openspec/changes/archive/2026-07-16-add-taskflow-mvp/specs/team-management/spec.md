## ADDED Requirements

### Requirement: 팀 생성
시스템은 인증된 사용자가 팀을 생성할 수 있도록 해야 한다(SHALL). 생성자는 자동으로 `owner_id`가 되고, 시스템은 `^[A-Z]{4}-[0-9]{4}$` 형식의 고유 초대코드를 자동 생성해야 하며(MUST), 생성자의 `users.team_id`를 즉시 갱신해야 한다.

#### Scenario: 정상 팀 생성
- **WHEN** 인증된 사용자가 팀 이름(1-30자)으로 `POST /teams`를 호출하면
- **THEN** 시스템은 `teams` 레코드를 생성하고 `201`과 함께 `id`, `name`, `invite_code`, `owner_id`를 반환하며 사용자의 `team_id`를 갱신한다

### Requirement: 초대코드로 팀 합류
시스템은 유효한 초대코드를 가진 사용자를 팀에 합류시켜야 한다(SHALL). 형식 오류, 존재하지 않는 코드, 이미 다른 팀 소속인 경우 각각 구분된 에러를 반환해야 한다(MUST).

#### Scenario: 정상 합류
- **WHEN** 팀 미가입 사용자가 유효한 초대코드로 `POST /teams/join`을 호출하면
- **THEN** 시스템은 `users.team_id`를 해당 팀으로 갱신하고 `200`과 팀 정보를 반환한다

#### Scenario: 초대코드 형식 오류
- **WHEN** 사용자가 `^[A-Z]{4}-[0-9]{4}$` 형식에 맞지 않는 코드로 합류를 시도하면
- **THEN** 시스템은 `400 { error: { code: 'VALIDATION_ERROR' } }`를 반환한다

#### Scenario: 존재하지 않는 초대코드
- **WHEN** 사용자가 존재하지 않는 초대코드로 합류를 시도하면
- **THEN** 시스템은 `404 { error: { code: 'NOT_FOUND' } }`를 반환한다

#### Scenario: 이미 다른 팀 소속
- **WHEN** 이미 `team_id`가 설정된 사용자가 다른 초대코드로 합류를 시도하면
- **THEN** 시스템은 `409 { error: { code: 'CONFLICT' } }`를 반환한다

### Requirement: 팀 정보 및 멤버 목록 조회
시스템은 팀 멤버에게 팀 정보와 멤버 목록(owner 여부 포함)을 제공해야 한다(SHALL). 비멤버의 조회 요청은 거부해야 한다(MUST).

#### Scenario: 멤버가 팀 정보 조회
- **WHEN** 팀 멤버가 `GET /teams/{id}`를 호출하면
- **THEN** 시스템은 `200`과 함께 팀 이름, 초대코드, member_count를 반환한다

#### Scenario: 멤버 목록 조회
- **WHEN** 팀 멤버가 `GET /teams/{id}/members`를 호출하면
- **THEN** 시스템은 각 멤버의 email, owner 여부, 합류일을 포함한 목록을 반환한다

#### Scenario: 비멤버 접근 차단
- **WHEN** 해당 팀 소속이 아닌 사용자가 `GET /teams/{id}` 또는 `GET /teams/{id}/members`를 호출하면
- **THEN** 시스템은 `403 { error: { code: 'FORBIDDEN' } }`를 반환한다

### Requirement: 팀 나가기
시스템은 사용자가 본인 계정으로 소속 팀을 나갈 수 있게 해야 한다(SHALL).

#### Scenario: 팀 나가기
- **WHEN** 팀 멤버가 `DELETE /teams/{id}/leave`를 호출하면
- **THEN** 시스템은 해당 사용자의 `users.team_id`를 NULL로 갱신하고 `200`을 반환한다

### Requirement: 팀 미가입 상태 강제 라우팅
시스템은 `team_id`가 NULL인 사용자가 칸반/채팅 등 팀 전용 화면에 접근하지 못하도록 막아야 한다(SHALL). 팀 소속 사용자가 다른 팀의 리소스에 URL로 직접 접근하는 경우도 동일하게 차단해야 한다.

#### Scenario: 미가입 사용자의 팀 리소스 접근
- **WHEN** `team_id`가 NULL인 사용자가 임의의 `/teams/{id}/*` 엔드포인트를 호출하면
- **THEN** 시스템은 `403 { error: { code: 'FORBIDDEN' } }`를 반환한다

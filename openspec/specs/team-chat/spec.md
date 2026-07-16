## ADDED Requirements

### Requirement: 메시지 전송
시스템은 팀 멤버가 팀 채팅에 메시지를 보낼 수 있게 해야 한다(SHALL). 메시지는 1000자를 초과할 수 없으며(MUST NOT), 클라이언트와 서버 양쪽에서 검증되어야 한다.

#### Scenario: 정상 전송
- **WHEN** 팀 멤버가 1000자 이내 내용으로 `POST /teams/{id}/messages`를 호출하면
- **THEN** 시스템은 `messages` 레코드를 생성하고 `201`을 반환한다

#### Scenario: 1000자 초과
- **WHEN** 팀 멤버가 1000자를 초과하는 내용으로 메시지 전송을 시도하면
- **THEN** 시스템은 `400 { error: { code: 'TOO_LONG', limit: 1000, actual: <실제길이> } }`를 반환한다

### Requirement: 폴링 기반 메시지 조회
시스템은 5초 간격 폴링을 지원하기 위해 `since=` 파라미터로 특정 시각 이후 메시지만 조회할 수 있어야 한다(SHALL). `since` 없이 조회하면 최근 메시지(최대 50개)를 반환해야 한다.

#### Scenario: 최초 진입 조회
- **WHEN** 팀 멤버가 `since` 없이 `GET /teams/{id}/messages`를 호출하면
- **THEN** 시스템은 해당 팀의 최근 메시지 목록(최대 50개, 발신자·시각 포함)을 반환한다

#### Scenario: 증분 폴링 조회
- **WHEN** 팀 멤버가 `GET /teams/{id}/messages?since=<마지막 메시지 시각>`을 호출하면
- **THEN** 시스템은 해당 시각 이후 생성된 메시지만 반환하며, 새 메시지가 없으면 빈 배열을 반환한다

### Requirement: 메시지 누락 없음 보장
시스템은 성공적으로 전송(201)된 메시지가 이후 모든 조회 요청에서 노출되도록 보장해야 한다(SHALL). 삭제된 메시지가 조회에서 빠지는 것은 누락으로 간주하지 않는다.

#### Scenario: 네트워크 재연결 후 동기화
- **WHEN** 클라이언트가 폴링 중 연결이 끊겼다가 복구되어 `since=` 파라미터로 재조회하면
- **THEN** 시스템은 끊긴 구간 동안 전송된 모든 메시지를 빠짐없이 반환한다

### Requirement: 메시지 삭제 권한
시스템은 메시지 작성자 본인만 자신의 메시지를 삭제할 수 있게 해야 한다(SHALL). 팀 owner라도 타인의 메시지는 삭제할 수 없다(MUST NOT).

#### Scenario: 본인 메시지 삭제
- **WHEN** 메시지 작성자가 `DELETE /messages/{id}`를 호출하면
- **THEN** 시스템은 메시지를 삭제하고 `200`을 반환한다

#### Scenario: 타인 메시지 삭제 시도 (owner 포함)
- **WHEN** 작성자가 아닌 사용자(팀 owner 포함)가 `DELETE /messages/{id}`를 호출하면
- **THEN** 시스템은 `403 { error: { code: 'NOT_OWNER' } }`를 반환한다

## ADDED Requirements

### Requirement: 태스크 생성
시스템은 팀 멤버가 팀 칸반에 태스크를 생성할 수 있게 해야 한다(SHALL). 생성 시 제목(1-100자)은 필수이며 상태는 기본값 `TODO`, `assignee_id`는 nullable로 지정 가능해야 한다.

#### Scenario: 정상 태스크 생성
- **WHEN** 팀 멤버가 제목과 선택적 담당자로 `POST /teams/{id}/tasks`를 호출하면
- **THEN** 시스템은 `tasks` 레코드를 `status=TODO`, `creator_id=현재사용자`로 생성하고 `201`을 반환한다

### Requirement: 칸반 조회 및 필터
시스템은 팀의 태스크 목록을 `TODO/DOING/DONE` 상태별로 조회할 수 있어야 하며(SHALL), 전체/`@me`(assignee 기준)/미할당 필터와 생성일 기준 정렬을 지원해야 한다.

#### Scenario: 전체 조회
- **WHEN** 팀 멤버가 `GET /teams/{id}/tasks`를 호출하면
- **THEN** 시스템은 해당 팀의 모든 태스크를 상태별로 반환한다

#### Scenario: 내 태스크 필터
- **WHEN** 팀 멤버가 `?filter=me`로 조회하면
- **THEN** 시스템은 `assignee_id = 현재사용자`인 태스크만 반환한다(`creator_id` 기준이 아님)

#### Scenario: 미할당 필터
- **WHEN** 팀 멤버가 `?filter=unassigned`로 조회하면
- **THEN** 시스템은 `assignee_id IS NULL`인 태스크만 반환한다

### Requirement: 태스크 상태 변경
시스템은 칸반 드래그로 인한 상태 변경을 `PATCH /tasks/{id}/status`로 전용 처리해야 한다(SHALL). 제목 수정과는 별도 API로 분리되어야 한다(MUST).

#### Scenario: 드래그로 상태 변경
- **WHEN** 팀 멤버가 카드를 다른 컬럼으로 드롭하여 `PATCH /tasks/{id}/status { status: 'DOING' }`를 호출하면
- **THEN** 시스템은 `tasks.status`를 갱신하고 `200`을 반환한다

### Requirement: 태스크 제목/담당자 수정
시스템은 태스크 상세 모달에서 제목과 담당자를 수정할 수 있게 해야 한다(SHALL).

#### Scenario: 제목 수정
- **WHEN** 팀 멤버가 `PUT /tasks/{id} { title }`를 호출하면
- **THEN** 시스템은 제목을 갱신하고 `200`을 반환한다(상태는 변경하지 않음)

### Requirement: 태스크 삭제 권한
시스템은 태스크 생성자(creator) 또는 팀 owner만 삭제를 허용해야 한다(SHALL). 그 외 사용자는 거부되어야 한다(MUST).

#### Scenario: 생성자 본인 삭제
- **WHEN** 태스크 생성자가 `DELETE /tasks/{id}`를 호출하면
- **THEN** 시스템은 태스크를 삭제하고 `200`을 반환한다

#### Scenario: owner의 타인 태스크 삭제
- **WHEN** 팀 owner가 다른 멤버가 생성한 태스크에 `DELETE /tasks/{id}`를 호출하면
- **THEN** 시스템은 태스크를 삭제하고 `200`을 반환한다

#### Scenario: 권한 없는 삭제 시도
- **WHEN** 생성자도 owner도 아닌 멤버가 `DELETE /tasks/{id}`를 호출하면
- **THEN** 시스템은 `403 { error: { code: 'FORBIDDEN' } }`를 반환한다

### Requirement: 태스크 단건 조회
시스템은 팀 멤버가 단일 태스크 상세를 조회할 수 있게 해야 한다(SHALL).

#### Scenario: 상세 조회
- **WHEN** 팀 멤버가 `GET /tasks/{id}`를 호출하면
- **THEN** 시스템은 제목, 상태, creator, assignee, 생성 시각을 반환한다

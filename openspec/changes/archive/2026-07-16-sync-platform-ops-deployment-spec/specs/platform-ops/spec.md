## MODIFIED Requirements

### Requirement: 환경 분리 (로컬/운영)
시스템은 `DATABASE_URL` 환경변수 하나로 SQLite와 Postgres 간 전환이 가능해야 한다(SHALL). 로컬 개발 및 Docker Compose 기본값은 이 PC에서 실행되는 로컬 Postgres 컨테이너이며, 필요하다면 SQLite로도 전환 가능하다. 외부 관리형 Postgres(Neon 등)는 선택 사항이며 기본 경로가 아니다. 로컬 SQLite 파일은 버전관리에서 제외되어야 한다.

#### Scenario: SQLite로 로컬 실행
- **WHEN** `DATABASE_URL`이 설정되지 않거나 SQLite 경로를 가리키면
- **THEN** 애플리케이션은 로컬 SQLite 파일을 사용해 기동한다

#### Scenario: 로컬 Postgres 컨테이너로 실행 (기본값)
- **WHEN** `DATABASE_URL`이 `docker-compose.yml`(또는 k8s/Helm)의 로컬 `postgres` 컨테이너를 가리키면
- **THEN** 애플리케이션은 동일한 코드로 해당 컨테이너에 연결하여 기동하며, 이 과정에 인터넷 연결이 필요하지 않다

#### Scenario: 외부 관리형 Postgres로 전환 (선택)
- **WHEN** `DATABASE_URL`이 외부 관리형 Postgres(Neon 등) 접속 문자열로 설정되면
- **THEN** 애플리케이션은 동일한 코드로 해당 Postgres에 연결하여 기동한다

### Requirement: 배포 파이프라인
시스템은 Docker Compose를 통한 로컬 배포(백엔드/프론트엔드/Postgres 모두 이 PC의 컨테이너)와, Kubernetes 매니페스트(`k8s/`) 또는 Helm 차트(`helm/taskflow`)를 통한 클러스터 배포를 지원해야 한다(SHALL). 두 경로 모두 기본값은 로컬에서 빌드한 이미지와 로컬(또는 클러스터 내부) Postgres이며, 외부 컨테이너 레지스트리나 클라우드 DB는 선택 사항이다.

#### Scenario: Docker Compose로 로컬 배포
- **WHEN** 개발자가 `docker compose up`(또는 `make docker-up`)을 실행하면
- **THEN** 백엔드, 프론트엔드, Postgres가 모두 이 PC의 컨테이너로 기동되고 서로 연결된다

#### Scenario: Kubernetes/Helm으로 클러스터 배포
- **WHEN** 개발자가 로컬에서 빌드한 이미지를 클러스터에 로드하고(`kind load docker-image` 등) `kubectl apply -f k8s/` 또는 `helm install`을 실행하면
- **THEN** 백엔드, 프론트엔드, (기본값인 경우) 클러스터 내부 Postgres가 기동되며 외부 레지스트리나 외부 DB 없이 동작한다

## ADDED Requirements

### Requirement: 정적 자산의 로컬 번들링 (외부 CDN 미사용)
시스템은 런타임에 필요한 정적 자산(프론트엔드 CSS 프레임워크, API 문서화 UI)을 외부 CDN에서 로드하지 않고 로컬에 번들링/내장해야 한다(SHALL). 빌드 시점의 일회성 도구 다운로드(예: Tailwind CLI, swagger-ui-dist 자산 벤더링)는 허용되나, 실행 중인 애플리케이션은 이 자산들을 위해 인터넷에 접근해서는 안 된다(MUST NOT).

#### Scenario: 프론트엔드 오프라인 렌더링
- **WHEN** 인터넷 연결 없이 프론트엔드를 로드하면
- **THEN** `frontend/tailwind.css`(로컬 빌드 산출물)로 스타일이 정상 렌더링된다

#### Scenario: API 문서 오프라인 접근
- **WHEN** 인터넷 연결 없이 `/docs`(Swagger UI)에 접속하면
- **THEN** `go:embed`로 내장된 swagger-ui 자산이 정상 로드되어 API 문서와 "Authorize" 기능이 동작한다

# 배포 절차 (Docker)

> Windows에 `make`가 없다면 각 타겟에 대응하는 명령을 `Makefile`에서 그대로 복사해 실행하면 됩니다 (예: `docker compose build`, `docker compose up -d --build`).

이 스택은 **기본값이 완전히 로컬 자족형(self-contained)** 입니다: DB(Postgres), 백엔드, 프론트엔드가 전부 이 PC의 Docker 컨테이너에서 뜨고, 외부 클라우드 서비스(Neon 등)나 CDN 없이 오프라인에서도 동작합니다. Tailwind CSS는 로컬 빌드 산출물(`frontend/tailwind.css`)로, Swagger UI 정적 자산은 백엔드 바이너리에 내장(`go:embed`)되어 있어 브라우저가 인터넷 CDN에 접근할 필요가 없습니다.

## 0. 준비

```bash
cp .env.example .env
# 필요시 .env의 JWT_SECRET 등을 수정 (DATABASE_URL 기본값은 아래 postgres 컨테이너를 가리킴)
```

## 1. Go 빌드만 확인하고 싶을 때 (Docker 없이)

```bash
make build   # backend/bin/server 생성
make test    # go test ./...
make run     # go run ./cmd/server (localhost:8000)
```

## 2. 프론트엔드 CSS 재빌드 (클래스 추가/변경 시에만)

```bash
make frontend-css
```

`index.html`/`js/*.js`에서 쓰는 Tailwind 클래스를 스캔해 `frontend/tailwind.css`를 다시 생성합니다. 최초 1회만 Tailwind standalone CLI를 받아오며(빌드 타임에만 인터넷 필요, `go mod download`와 동일한 성격), 실행 중인 앱은 이 정적 파일만 읽으므로 런타임에는 인터넷이 필요 없습니다.

## 3. Docker 이미지 빌드

```bash
make docker-build
# 또는: docker compose build
```

- `backend/Dockerfile`: 멀티스테이지 빌드. `golang:1.23-alpine`에서 `CGO_ENABLED=0`으로 정적 바이너리를 빌드하고, 최종 이미지는 `alpine:3.20`에 바이너리 + `docs/`(Swagger 스펙) + 내장된 swagger-ui 자산만 담아 가볍게 유지합니다. SQLite 드라이버가 순수 Go 구현(`modernc.org/sqlite`)이라 cgo/gcc 없이도 빌드됩니다.
- `frontend/Dockerfile`: 정적 파일(빌드된 `tailwind.css` 포함)을 `nginx:1.27-alpine`으로 서빙합니다.

## 4. 로컬 기동 (Postgres 컨테이너가 기본값)

```bash
make docker-up
```

- 백엔드: http://localhost:8000 (Swagger: http://localhost:8000/docs — 오프라인에서도 정상 동작)
- 프론트엔드: http://localhost:3000
- DB: `postgres` 컨테이너(이 PC 안에서 실행), 데이터는 `postgres-data` 볼륨에 영속화
- 백엔드는 `postgres` 컨테이너가 healthy 상태가 될 때까지 기다렸다 기동합니다(`depends_on: condition: service_healthy`).

로그 확인:

```bash
make docker-logs
```

종료:

```bash
make docker-down
```

### SQLite로 바꾸고 싶다면 (postgres 컨테이너 없이 backend 단독)

`.env`에서 `DATABASE_URL=sqlite:///./data/taskflow.db`로 바꾸면 됩니다. 다만 `docker-compose.yml`의 `depends_on: postgres`는 그대로 남아있으므로, postgres 컨테이너 자체는 계속 함께 뜹니다(단순 미사용 상태). 완전히 분리하고 싶다면 `docker-compose.yml`에서 backend의 `depends_on` 블록과 `postgres` 서비스를 직접 지워도 됩니다.

## 5. 외부 클라우드 DB(Neon 등)로 전환하고 싶을 때 (선택 사항, 기본값 아님)

기본값은 위 4번의 로컬 `postgres` 컨테이너입니다. 그래도 클라우드 Postgres(Neon 등)를 쓰고 싶다면, `.env`(또는 실제 배포 환경의 시크릿 관리 도구)의 `DATABASE_URL`만 교체하면 됩니다. 코드 변경은 필요 없습니다 — `backend/internal/db/db.go`가 `DATABASE_URL` 스킴(`sqlite://` vs `postgres://`)에 따라 자동으로 분기합니다.

```
DATABASE_URL=postgres://<user>:<password>@<neon-host>/<db>?sslmode=require
```

`JWT_SECRET`은 운영에서 반드시 강력한 랜덤 값으로 교체하고, `CORS_ORIGINS`는 실제 프론트엔드 배포 도메인으로 좁혀야 합니다.

## 6. 참고 — 이 환경에서 검증한 것과 검증하지 못한 것

- ✅ `CGO_ENABLED=0 GOOS=linux go build`로 리눅스 대상 정적 바이너리 빌드 성공 확인
- ✅ `docker compose config`로 `docker-compose.yml` 문법 및 환경변수 치환 정상 확인
- ✅ Tailwind CSS 로컬 빌드(`frontend/tailwind.css`) 및 index.html에서 CDN 제거 확인
- ✅ Swagger UI 자산(`swagger-ui.css`, `swagger-ui-bundle.js`) `go:embed`로 내장 후 `/docs`, `/docs/swaggerui/*` 로컬 서빙 확인 (`go test`, 로컬 서버 curl로 200 확인)
- ⚠️ 이 환경에는 Docker 데몬(Docker Desktop 엔진)이 떠 있지 않아 실제 `docker compose build`/`up`(postgres 컨테이너 기동 포함)까지는 확인하지 못했습니다. 로컬에서 Docker Desktop을 켠 상태로 `make docker-up`을 실행해 최종 확인해주세요.

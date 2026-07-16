# 배포 절차 (Docker)

> Windows에 `make`가 없다면 각 타겟에 대응하는 명령을 `Makefile`에서 그대로 복사해 실행하면 됩니다 (예: `docker compose build`, `docker compose up -d --build`).

## 0. 준비

```bash
cp .env.example .env
# 필요시 .env의 JWT_SECRET, DATABASE_URL 등을 수정
```

## 1. Go 빌드만 확인하고 싶을 때 (Docker 없이)

```bash
make build   # backend/bin/server 생성
make test    # go test ./...
make run     # go run ./cmd/server (localhost:8000)
```

## 2. Docker 이미지 빌드

```bash
make docker-build
# 또는: docker compose build
```

- `backend/Dockerfile`: 멀티스테이지 빌드. `golang:1.23-alpine`에서 `CGO_ENABLED=0`으로 정적 바이너리를 빌드하고, 최종 이미지는 `alpine:3.20`에 바이너리 + `docs/`(Swagger 스펙)만 담아 가볍게 유지합니다. SQLite 드라이버가 순수 Go 구현(`modernc.org/sqlite`)이라 cgo/gcc 없이도 빌드됩니다.
- `frontend/Dockerfile`: 정적 파일을 `nginx:1.27-alpine`으로 서빙합니다.

## 3. 로컬 기동 (SQLite, 기본값)

```bash
make docker-up
```

- 백엔드: http://localhost:8000 (Swagger: http://localhost:8000/docs)
- 프론트엔드: http://localhost:3000
- SQLite 파일은 `backend-data` 볼륨에 영속화되어 컨테이너를 내려도 유지됩니다.

로그 확인:

```bash
make docker-logs
```

종료:

```bash
make docker-down
```

## 4. 로컬 Postgres로 테스트하고 싶을 때

운영(Neon Postgres)에 배포하기 전, 동일한 Postgres 방언으로 먼저 확인하고 싶다면:

```bash
# .env에서 DATABASE_URL을 아래로 변경
# DATABASE_URL=postgres://taskflow:taskflow@postgres:5432/taskflow?sslmode=disable

make docker-up-postgres
```

`docker-compose.yml`의 `postgres` 서비스는 `profiles: ["postgres"]`로 묶여있어 기본 `docker compose up`에는 포함되지 않고, `--profile postgres`를 붙였을 때만 함께 뜹니다.

## 5. 운영(Neon) 배포로 전환

운영에서는 `docker-compose.yml`의 `postgres` 서비스 대신, `.env`(또는 실제 배포 환경의 시크릿 관리 도구)의 `DATABASE_URL`을 Neon 접속 문자열로 교체하면 됩니다. 코드 변경은 필요 없습니다 — `backend/internal/db/db.go`가 `DATABASE_URL` 스킴(`sqlite://` vs `postgres://`)에 따라 자동으로 분기합니다.

```
DATABASE_URL=postgres://<user>:<password>@<neon-host>/<db>?sslmode=require
```

`JWT_SECRET`은 운영에서 반드시 강력한 랜덤 값으로 교체하고, `CORS_ORIGINS`는 실제 프론트엔드 배포 도메인으로 좁혀야 합니다.

## 6. 참고 — CI/이 환경에서 검증한 것과 검증하지 못한 것

- ✅ `CGO_ENABLED=0 GOOS=linux go build`로 리눅스 대상 정적 바이너리 빌드 성공 확인
- ✅ `docker compose config`로 `docker-compose.yml` 문법 및 환경변수 치환 정상 확인
- ⚠️ 이 환경에는 Docker 데몬(Docker Desktop 엔진)이 떠 있지 않아 실제 `docker compose build`/`up` 실행까지는 확인하지 못했습니다. 로컬에서 Docker Desktop을 켠 상태로 `make docker-up`을 실행해 최종 확인해주세요.

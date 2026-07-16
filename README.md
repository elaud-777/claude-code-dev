# TaskFlow MVP

소규모 팀(3-5인)이 칸반 보드와 팀 채팅을 한 화면에서 오가며 업무 진행을 추적하는 웹앱입니다. 인증(JWT), 팀/초대코드, 칸반(TODO/DOING/DONE), 팀 채팅(폴링) 기능을 제공합니다.

- **백엔드**: Go (chi + sqlx + golang-jwt + bcrypt), Swagger UI 내장
- **프론트엔드**: Vanilla JS + Tailwind CSS (로컬 빌드, CDN 미사용)
- **DB**: Postgres(로컬 Docker 컨테이너가 기본값) 또는 SQLite
- 전체 스택이 **외부 클라우드 서비스 없이 이 PC의 Docker 컨테이너만으로** 동작하도록 구성되어 있습니다.

## 디렉토리 구조

```
backend/    Go 백엔드 (API, DB, 테스트) — backend/README.md
frontend/   Vanilla JS + Tailwind 프론트엔드
openspec/   기능 명세, 설계 문서, 변경 이력 (OpenSpec)
docker-compose.yml  로컬 개발/실행 (백엔드+프론트엔드+Postgres 전부 컨테이너)
k8s/        플레인 Kubernetes 매니페스트 — k8s/README.md
helm/       Helm 차트 (templates 기반 배포) — helm/README.md
DEPLOY.md   Docker 기반 배포 절차 상세
```

## 빠른 시작 (로컬, Docker)

```bash
cp .env.example .env
make docker-up
```

- 프론트엔드: http://localhost:3000
- 백엔드: http://localhost:8000 (API 문서: http://localhost:8000/docs)
- DB: 로컬 `postgres` 컨테이너 (docker-compose에 포함, 별도 설치 불필요)

자세한 절차와 옵션(로컬 Go 빌드만 하기, SQLite로 전환하기, 외부 Postgres로 전환하기 등)은 [DEPLOY.md](DEPLOY.md)를 참고하세요.

## Kubernetes 배포

- 플레인 매니페스트: [k8s/README.md](k8s/README.md)
- Helm 차트: [helm/README.md](helm/README.md)

둘 다 기본값이 로컬 빌드 이미지 + 클러스터 내부 Postgres로, 별도 컨테이너 레지스트리나 외부 DB 없이 로컬 클러스터(kind/minikube)에서 바로 띄울 수 있습니다.

## 개발 히스토리 (OpenSpec)

이 프로젝트는 [OpenSpec](openspec/) 워크플로우(propose → apply → archive)로 진행되었습니다. 최초 Python(FastAPI) 구현 후 Go로 전면 재작성했고, 관련 제안/설계/스펙은 `openspec/changes/archive/`에 남아 있습니다. 현재 확정된 기능 스펙은 `openspec/specs/`를 참고하세요.

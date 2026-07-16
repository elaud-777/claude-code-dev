## 1. 스펙-구현 일치 확인 (코드 변경 없음, 검증만)

- [x] 1.1 `docker-compose.yml`의 `postgres` 서비스가 기본(비-optional) 서비스이고 `DATABASE_URL` 기본값이 이를 가리키는지 재확인
- [x] 1.2 `k8s/postgres.yaml`, `helm/taskflow/templates/postgres-*.yaml`이 기본값으로 존재하고 `helm template` 렌더링이 정상인지 재확인
- [x] 1.3 `frontend/index.html`에 `cdn.tailwindcss.com` 참조가 없고 `frontend/tailwind.css`를 로드하는지 재확인
- [x] 1.4 `backend/cmd/server/swagger.go`가 `unpkg.com`이 아닌 `go:embed` 자산을 서빙하는지 재확인
- [x] 1.5 `openspec/specs/platform-ops/spec.md`에 더 이상 Vercel/Neon-필수 문구가 남아있지 않은지 확인 (아카이브 시 델타 적용 후)

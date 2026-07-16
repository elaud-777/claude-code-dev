# 이 프로젝트에서 작업할 때 지켜야 할 규칙

이 문서는 저장소에 커밋되어 어느 PC에서 클론하든 동일하게 적용됩니다. (로컬 메모리 `~/.claude/...`는 PC/계정별로 저장되어 다른 PC로 옮겨지지 않으므로, 이 프로젝트에 반드시 적용되어야 하는 규칙은 여기 문서화합니다.)

## 1. 항시 OpenSpec 절차를 지킬 것

애플리케이션 코드든, 인프라/배포 도구(Docker, Kubernetes, Helm, CI 등)든 상관없이 사소한 한 줄 수정이나 순수 문서 편집이 아닌 이상 **반드시 OpenSpec 절차**(`/opsx:propose` → `/opsx:apply` → `/opsx:archive`)를 거쳐서 진행합니다. 먼저 구현하고 나중에 스펙에 동기화하는 방식은 지양합니다.

**이유**: 과거 이 프로젝트에서 Docker Compose, Kubernetes 매니페스트, Helm 차트, CDN 제거 작업을 OpenSpec 절차 없이 직접 구현했다가, `openspec/specs/platform-ops/spec.md`가 실제 시스템(Vercel/Neon 언급)과 어긋나는 문제가 발생했습니다. 이를 뒤늦게 발견해 `sync-platform-ops-deployment-spec`라는 별도 change로 동기화해야 했습니다. "인프라/도구성 작업이라 스펙과 무관하다"는 판단은 하지 않습니다 — 배포 방식이나 운영 방식이 바뀌면 그것도 스펙 대상입니다.

## 2. 커밋 후 자동으로 푸시할 것

이 저장소(`elaud-777/claude-code-dev`)에서 git commit을 만들면, 별도로 "푸시해줘"라고 확인받지 않고 바로 `git push origin master`까지 진행합니다. 단, 무엇을 커밋/푸시했는지는 항상 결과로 보여줍니다.

## 참고

- 명세/설계/변경 이력: `openspec/specs/`(현재 확정 스펙), `openspec/changes/archive/`(완료된 change 이력)
- 새 작업을 시작하려면: `/opsx:propose "<하고 싶은 것 설명>"`

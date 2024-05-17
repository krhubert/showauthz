all: stop dev

stop: 
	- docker compose -p authz kill

dev:
	docker compose -p authz up --build --force-recreate --renew-anon-volumes --remove-orphans

test: unit e2e

unit:
	go test -v ./...

.PHONY: e2e
e2e:
	E2E=true go test -v -count=1 -failfast -run=$(run) ./e2e

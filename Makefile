.PHONY: run

run:
	echo "Run using room/"

test: test_integration

test_unit:
	./run-tests.sh

test_integration:
	docker compose --env-file default.env -f compose.integration.yml up --build --abort-on-container-exit --exit-code-from test-runner
	docker compose --env-file default.env -f compose.integration.yml down

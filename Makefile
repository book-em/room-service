.PHONY: run

run:
	echo "Run using room/"

test: test_integration

test_unit:
	./run-tests.sh


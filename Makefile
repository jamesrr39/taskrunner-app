
.PHONY: install_dependencies
install_dependencies:
	apt install build-essential libgtk2.0-dev

.PHONY: run
run:
	go run taskrunner-app-main.go

.PHONY: install
install:
	test -n "${GOBIN}"
	go build -o ${GOBIN}/taskrunner-app taskrunner-app-main.go

.PHONY: test
test:
	go test ./...
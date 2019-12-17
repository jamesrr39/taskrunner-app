
.PHONY: install_dependencies
install_dependencies:
	apt install build-essential libgtk2.0-dev

.PHONY: run
run:
	go run taskrunner-app-main.go

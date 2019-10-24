.PHONY: run
run:
	source .env && \
	go run main.go

.PHONY: build
build:
	docker build . -t argocd-slack-notifier

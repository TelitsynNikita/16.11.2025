.PHONY: run-service
run-service: ## запускаем сервис через go run .
	go run cmd/main.go

.PHONY: build-service
build-service: ## билдим бинарник
	go build -o app cmd/main.go

.PHONY: build-and-run-service-by-unix
build-and-run-service: ## билдим бинарник и запускаем его как процесс Unix ОС
	go build -o app cmd/main.go
	./app

.PHONY: build-and-run-service-by-windows
build-and-run-service: ## билдим бинарник и запускаем его как процесс Windows ОС
	go build -o app cmd/main.go
	./app.exe
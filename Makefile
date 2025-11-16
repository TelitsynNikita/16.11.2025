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

send-urls: ## Отправить сообщение в link/check_by_urls
	 curl --location 'http://localhost:8080/link/check_by_urls' \
     --header 'Content-Type: application/json' \
     --data '{"links": ["google.com", "ya.ru", "https://blog.ildarkarymov.ru/posts/graceful-shutdown/", "https://purpleschool.ru/knowledge-base/article/os-file"]}'

send-ids: ## Отправить сообщение в link/check_by_id
	 curl --location 'http://localhost:8080/link/check_by_id' \
     --header 'Content-Type: application/json' \
     --data '{"links_list": [1]}'
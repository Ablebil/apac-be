run:
	@docker compose restart
build:
	@docker compose down
	@docker compose up --build
stop:
	@docker compose stop
remove:
	@docker compose down
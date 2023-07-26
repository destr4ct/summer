.PHONY: run, logs

buildPath = builds

bpc:
	 go build -o $(buildPath)/producer destr4ct/summer/cmd/producer
	 go build -o $(buildPath)/consumer destr4ct/summer/cmd/consumer

tg:
	go build -o $(buildPath)/tg destr4ct/summer/cmd/telegram

up:
	docker compose up --build --remove-orphans -d
logs:
	docker compose logs
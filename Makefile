up:
	docker compose down
	docker compose up --build -d
	docker image prune -f

ps:
	docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"

log:
	git log --oneline -n 10
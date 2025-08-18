deps:
	chmod +x setup_docker.sh
	./setup_docker.sh
	docker --version
	docker compose version

local:
	colima start
	docker build -t huntergoff .
	docker run --env-file .env --rm -p 8080:8080 huntergoff

remote:
	fly deploy
	open https://huntergoff.com

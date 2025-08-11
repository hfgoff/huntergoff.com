local:
	colima start
	docker build -t huntergoff .
	docker run -p 8080:8080 huntergoff

remote:
	fly deploy
	open https://huntergoff.com
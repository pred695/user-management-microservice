build:
	sudo chmod 777 -R pgdata
	docker build -t user-management-microservice .
compose: 
	docker-compose up

decompose:
	docker-compose down
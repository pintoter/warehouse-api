run:
	docker-compose -f docker-compose.yml up --remove-orphans

stop:
	docker-compose down --remove-orphans

lint:
	make lint
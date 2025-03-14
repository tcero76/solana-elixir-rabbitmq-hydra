QUEUE_NAME=solana_event_queue
MESSAGE={"event": "test_event", "data": "Hello, Solana!"}

.PHONY: sendmsg setupRabbitmq up kill down watch build exec

watch:
	watch docker ps -a

up:
	docker compose up -d

kill:
	docker rm -f $(filter-out $@,$(MAKECMDGOALS))

down:
	docker compose down

build:
	docker compose build $(filter-out $@,$(MAKECMDGOALS)) --no-cache --progress=plain

exec:
	docker exec -it $(filter-out $@,$(MAKECMDGOALS)) /bin/bash

sendmsg:
	docker exec broker rabbitmqadmin publish \
		routing_key=$(QUEUE_NAME) \
		payload='$(MESSAGE)'

setupRabbitmq:
	docker exec $(RABBITMQ_CONTAINER) rabbitmqadmin declare queue name=$(QUEUE_NAME) durable=true

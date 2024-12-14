.PHONY: lint

lint:
	golangci-lint run -c .golangci.yaml

run: .create-secret .build .deploy

stop:
	docker stack rm auth

# смотри readme
.create-secret: # AHTUNG! так нельзя делать. но для упрощения проверки тестового можно :D
	@if ! docker secret ls | grep -q jwt; then \
		echo "Creating JWT secret..."; \
		echo "SPY" > ./jwt_secret.txt; \
		docker secret create jwt ./jwt_secret.txt; \
		rm ./jwt_secret.txt; \
	else \
		echo "JWT secret already exists."; \
	fi

.build:
	docker build -t auth_app .

.deploy:
	docker stack deploy -c docker-compose.yml auth


# это для тестов без секрета

compose-run: .compose-build .compose-up

.compose-build:
	docker compose build

.compose-up:
	docker compose up -d

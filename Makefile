.PHONY: lint

lint:
	golangci-lint run -c .golangci.yaml

run: .create-secret .build .test .deploy

stop:
	docker stack rm auth

.create-secret: # AHTUNG! так нельзя делать. но для упрощения проверки тестового можно :D
		@if not exist jwt_secret.txt ( \
    		echo Checking if JWT secret exists in Docker... & \
    		docker secret ls | findstr jwt > nul && ( \
    			echo JWT secret already exists. \
    		) || ( \
    			echo Creating JWT secret... & \
    			echo SPY > jwt_secret.txt & \
    			docker secret create jwt jwt_secret.txt & \
    			del jwt_secret.txt \
    		) \
    	) else ( \
    		echo JWT secret file already exists locally. \
    	)

.build:
	docker build -t auth_app .

.test:
	go test ./...

.deploy:
	docker stack deploy -c docker-compose.yml auth


# это для тестов без секрета

compose-run: .compose-build .compose-up

.compose-build:
	docker compose build

.compose-up:
	docker compose up -d

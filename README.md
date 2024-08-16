**This repository contains the design and implementation of a text streaming api.**

**ENTITIES**
1. There are 3 workers inside the cmd folders (`tgione, tgitwo, tgithree`) which support one endpoint each (`/text/generate`) which mocks the characteristics of a text generation interface api. Also, each worker supports a  `/ping` endpoint for the health check of the services.
2. There is a `redis` instance to work as a db, which stores the health check of each of the services, if any of the text generation services goes down, the text streaming api should not make a request to that service.
3. There is another api, `text/health/check` , which should be put on the cloud schedulers, which checks the status of each of the tgi services and mark their response status in the redis.

**COMMAND TO RUN THE APPLICATION**
1. There is a `docker-compose.yaml` file in the repo to run the latest redis image. `docker compose up -d`
2. Start the main application server. `make start_server` / `go run cmd/app/*.go`
3. Start the tgione worker. `make start_tgione` / `go run cmd/tgione/main.go`
4. Start the tgitwo worker. `start_tgitwo` / `go run cmd/tgitwo/main.go`
5. Start the tgithree worker. `start_tgithree` / `go run cmd/tgithree/main.go`
6. Import the postman collection attached to make api calls.

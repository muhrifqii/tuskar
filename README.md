# Tuskar, a Tasking System Foo Bar

## Required Tools to Run
- Docker & docker-compose
- CMake

## Getting Started
### Env Setup
First thing first after all tools installed, prepare the environment variables
```bash
# copy the example.env to .env
cp .sample.env .env
```

And continued with preparing the containers (only once)
```bash
make dev-env
```
### Database Setup
After all containers are ready, you can run the migration script using the following command
```bash
make migrate-up
```
It will install golang-migrate to `./bin` directory. And then when prompt showed up asking how many steps you want to perform migration, just press enter or put 1
```bash
How many migration you wants to perform (default value: [all]): [Your input goes here]
```
When the migration is done, you can run the application using air. But if something goes wrong and getting dirty, you need to run the following command to clean up the database, followed by `1`
```bash
# only run if migration schema getting dirty
make migrate-dirty
# then run migrate-down and press Enter or put 1
make migrate-down
```

## Run The Application
You can check all the port config in the `.sample.env` file, default server port is `8081`
```bash

# Run the application
make up

```
### Health Check
Health check is included so you can check the liveness and readiness by hitting the following endpoints
```bash
curl --location --request GET 'localhost:8081/health'
curl --location --request GET 'localhost:8081/ready'
```

### Authenticate
The Task API are authn protected, you need to login first before accessing the API
```bash
curl --location --request POST 'localhost:8081/authenticate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "saitama",
    "password": "saitama"
}'
```

## Logs
Promtail + Loki are used and the logs are stored in the `infras/logs` directory. It is also possible to use Grafana to visualize the logs. Grafana is accessible on `localhost:3030`

# VB-Events

## How to run locally

There is a .env file in the backend directory that you can use to set the environment variables.
Defaults for local development are

```bash
DBHost      localhost
DBPort      5432
DBUser      vbevents
DBPassword  vbevents
DBDatabase  events
SSLMode     disabled
APIHost     localhost
APIPort     8200
CORSOrigins http://localhost:5173
JWTSecret   (unset by default)
```

Most of these values have a fallback value in the code. However, the JWT secret is required and must be set in the .env file or the environment.
You can generate a JWT secret with the following command

```bash
openssl rand -base64 32
```

For the DBHost, DBPort, DBUser, and DBPassword, you can look in the docker-compose.yml or the /testing/start-local-db.sh script.
For a live environment, you will need to use a secrets manager to store the user and password for deployment.

For SSLMode, because this is demo code, it is set to disabled.
https://www.postgresql.org/docs/current/libpq-ssl.html

For the APIHost, I recomend keeping it localhost if you are running the backend locally or on 0.0.0.0 if you are running it in Docker.
For the APIPort, I suggest just using the default of 8200.

For the CORS Origins, you can use the following command to get the origin of the frontend
The fallback is set to the Vite Dev Server default.
In docker-compose.yml this is set to http://localhost:8080

### Without Docker

Ensure you have a Postgres database running that you can connect to and update the .env file

If you don't have postgres you can use 
```bash
/testing/start-local-db.sh
```
this will start a postgres database on port 5432 and create the database and user.
there is a default user and password for the database for local development.

#### How to run the backend
To run the backend ensure you have Go installed. The backend can be run locally without Docker.

```bash
cd backend
go run cmd/main.go
## or you can use the Makefile
make run
```

If you want to run the backend with hot reloading you can run
```bash
air
```

Other commands that you can run are 
```bash
make docs # Generates Swagger docs
make build # Builds the backend
make run # Runs the backend
make fmt # Formats the code
make tidy # go mod tidy for dependencies
make localdb # Starts a local Postgres database
make test # Runs the tests
make seed # Seeds the database with demo data - only run in development
make api-test # Runs the Postman tests
```

#### How to run the frontend
Node 24.37.0 is required to run the frontend.

```bash
cd frontend
npm ci
npm run dev

```
This will install the dependencies and run the frontend in development mode.
Other commands that you can run are
```bash
npm run build # Building the frontend
npm run type-check # Type checking the frontend
npm run lint  
npm run lint:oxlint
npm run lint:eslint
npm run format
```

I have not added any tests to the frontend yet.

### With Docker

To run the backend with Docker ensure you have Docker and Docker Compose installed.

```bash
# From the project root
docker compose up
```

Please read the docker-compose.yml file details on environment variables set and how the healthcheck works.

A note on the two sets of database variables you will see:
`DB_DATABASE`, `DB_USER` etc are what the Go backend reads to connect to the database (see backend/internal/config).
`POSTGRES_DB`, `POSTGRES_USER` etc are only used by docker-compose to initialise the Postgres container on first start.
In the docker-compose.yml the backend's `DB_*` values are filled in from the `POSTGRES_*` values so there is a single source of truth and they cannot drift apart.
You only need to set the `DB_*` variables yourself when running the backend without Docker.


## Architectural Decisions

I chose Go because its lightweight and easy to use. I chose Vue because its ideal for this sort of project.

To keep things simple, it can be easier to serve the frontend from the backend. However, because this may be hosted on Kubernetes, its easier to serve the frontend from a separate container.
I could have used an auto-incrementing ID for the events and the participants or even a UUID v4. 
Auto-incrementing IDs are not ideal because they can be easily guessed. 
UUID v4 is not ideal its harder for the database to index and query.
However, UUID v7 is time-based and much faster for the database to work with.

Instead of keeping a list of participants per event or a list of events per participant, I used a many-to-many relationship between the two.
This means if you want to collect all of the events a participant has joined, you can do that with a single query.
Or if you want to collect all of the participants who have joined a specific event, you can do that with a single query.
The join-table is called `participant_event` and it has a `participant_id` and an `event_id` column. Making it smaller and faster to query using the UUID v7 IDs.

### Structure of the code

The backend has a cmd directory that contains the main.go file and the seed command that is used to seed the database for demos.
The internal directory contains the config, db, handlers, middleware, models, repository, and services directories.
config contains the Load function that is used to load the configuration from the .env file or the environment.
db contains the database connection and the initDB function that is used to connect to the database.
handlers contains the code used to handle the HTTP requests.
middleware is used for the authentication middleware and any other middleware that is added.
models contains the models used by the application. These are GORM structs because I felt it would be easier for this use-case.
repository contains the code used to interact with the database.
services contains the code used to interact with the services.

The frontend structure is a bit different.
The models are in the src/types directory. 
The services are in the src/services directory and are used to interact with the backend.
The store is mostly for keeping track of the authentication state.
The components are in the src/components directory and are used to display the Dialogs UI similar 
to the views directory which is used for the pages.


## CI/CD
See the .github/workflows directory

## Prometheus Metrics
The backend exposes prometheus metrics on /metrics
Remember to use the backend port, not the frontend port. `localhost:8200/metrics`

## Postman Tests
Postmans are in the testing/postman directory.
The tests are run with the `make api-test` command.
Tests are also run in the CI/CD pipeline.

## Docs
The backend exposes Swagger documentation on /swagger
Remember to use the backend port, not the frontend port. `localhost:8200/swagger`

## Logs
The backend logs are in a JSON format and are output to stdout.
This should be ideal for logging in a containerized environment.
However, there is a RequestLogger middleware used that can be noisey.
You may want remove from main.go for sending logs to Datadog or Splunk as it can be expensive.

## Dependencies Security
I recommend installing trivy locally to scan the images for vulnerabilities.
Build the images with the following command

```bash
docker build ./backend --file Dockerfile --tag vb-events-backend:latest
docker build ./frontend --file Dockerfile --tag vb-events-frontend:latest
```

Then run the following command to scan the images for vulnerabilities.

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image --exit-code 1 --severity HIGH,CRITICAL vb-events-backend:latest
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy image --exit-code 1 --severity HIGH,CRITICAL vb-events-frontend:latest
```

## With more time
If it was a microservice project you would need tracing added.

A location to upload the docker images to.

K6 load testing

More unit tests and integration tests.

Permissions for the users to be able to be limited.
- Only allowing them to remove their own events
- Preventing them from removing other people's events
- Inviting other people wouldn't automatically add them to the event
- etc

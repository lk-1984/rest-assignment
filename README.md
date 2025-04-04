# rest-assignment

The `rest-assignment` is an example project of how to create a REST API with `Go` language, `Gin Web Framework`, and `PostgreSQL`. The API application is simple Go based REST API that manages (CRUD) three resources: continents, countries, and cities. It is written with Macbook that has MacOS Sequioa 15.3.2, but it should work in Linux too.

There is a devcontainer for `Visual Studio Code` in which the Go application can be developed. The devcontainer has Go, and it is connected to the Docker compose network, in which the application is running, after it has been deployed. That allows to run the compiled Go code against the same PostgreSQL instance, as the containerized Go application.

## Run Devcontainer

Clone the Git repository, and go into the root directory of the repository.

Create Docker network by running `docker network create api-network`.

Type `ctrl + shift + p` to starting the devcontainer Linux, or `shift + command + p` in Macbook to start the devcontainer.

Once it is running, type `make` to list all the make targets.

## Build API application

Run `make build`. The executable `api` is located under `./build/`.

## Run unit tests of API Application

Run `make test/unit`.

## Run API application with Postgres

Open a terminal from your host computer. Go into the directory of the Git repository.

Run `make docker/compose/postgres/up`.

Once PostgreSQL is running, then go back into the devcontainer, and run `make run`, and the API app starts, and it is connected into the PostgreSQL database.

To stop, press `ctrl + c` to stop the API app on the devcontainer, and then go into the host terminal, and run `make docker/compose/postgres/down`.

## Run API application in Docker compose

Run the API app and PostgreSQL in Docker compose by running `make docker/compose/up`.

## Usage

Run `export BASE_URL="http://localhost:8080/api/v1"`, or `export BASE_URL="http://api:8080/api/v1"` depending on where you run the curl commands - at host, or at devcontainer by opening another terminal in the Visual Studio Code. In this example, the latter one is used, and the API app, and PostgreSQL, are assumed to be running in Docker Compose.

Get all continents. The response should be an empty set, if you have started this for the first time.

```
curl -X GET ${BASE_URL}/continents
```

Create a continent.

```
curl -X POST -H "Content-Type: application/json" -d '{"name":"Asia"}' ${BASE_URL}/continent
```

Get the continent id from the response, and get the continent with that id.

```
curl -X GET ${BASE_URL}/continent/<id>
```

We made a mistake, let's change the Asia into Europe.

```
curl -X PUT -H "Content-Type: application/json" -d '{"name":"Europe"}' ${BASE_URL}/continent/<id>
```

Check that the content has been updated, and use the same id as you did before. You should see Europe instead of Asia.

```
curl -X GET ${BASE_URL}/continent/<id>
```

Now, let's add country to the continent.

```
curl -X POST -H "Content-Type: application/json" -d '{"name":"Germany","continent_id":<id>}' ${BASE_URL}/country
```

Get the id of the country, and get the country by the id.

```
curl -X GET ${BASE_URL}/country/<id>
```

Let's add another country.

```
curl -X POST -H "Content-Type: application/json" -d '{"name":"Poland","continent_id":<id>}' ${BASE_URL}/country
```

Let's view all the countries in the continent that was created earlier.


```
curl -X GET ${BASE_URL}/countries
```

Oops, we made a mistake, again. Let's fix it. Take the id of the first country.

```
curl -X PUT -H "Content-Type: application/json" -d '{"name":"France","continent_id":<id>}' ${BASE_URL}/country/<id>
```

Now, view all countries again, and you should see France and Poland.

```
curl -X GET ${BASE_URL}/countries
```

Then, add a city into France.

```
curl -X POST -H "Content-Type: application/json" -d '{"name":"Nice","country_id":<id>}' ${BASE_URL}/city
```

Add another city, but into Poland.

```
curl -X POST -H "Content-Type: application/json" -d '{"name":"Warsaw","country_id":<id>}' ${BASE_URL}/city
```

List all cities.

```
curl -X GET ${BASE_URL}/cities
```

Oh, we made another mistake. Let's fix it. Use the id of the first country.

```
curl -X PUT -H "Content-Type: application/json" -d '{"name":"Paris","country_id":<id>}' ${BASE_URL}/city/<id>
```

Check the cities again, and you should see Warsaw and Paris.

```
curl -X GET ${BASE_URL}/cities
```

Then, let's wrap up, and delete everything, starting from cities. Run this for each of the city ids.

```
curl -X DELETE ${BASE_URL}/city/<id>
```

Check that there are no cities left.

```
curl -X GET ${BASE_URL}/cities
```

Then delete countries.

```
curl -X DELETE ${BASE_URL}/country/<id>
```

Check that there are no countries left.

```
curl -X GET ${BASE_URL}/countries
```

Then delete the continent.

```
curl -X DELETE ${BASE_URL}/continent/<id>
```

Verify that there isn't any continents left.

```
curl -X GET ${BASE_URL}/continents
```

That's it, now you have succesfully tested the API application in action.
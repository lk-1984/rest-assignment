# API

1. Pagination for responses having over x items.
2. Improve validation of requests and responses.
3. Filtering with URL query parameters.
4. Error masking so that internal problems are not exposed.
5. Keycloak integration (OpenID Connect, OAuth) for authentication and authorization.
6. Prometheus metrics for requests, responses, and Go. PostgreSQL also.
7. Grafana dashboards for the metrics.
8. OpenAPI specification.
9. Log formatting, JSON for instance.

# Database

1. Create ER diagram.
2. Disallow duplicates on continent, country, and city names.
3. Backup and restore.

# Unit Tests

1. Add OK cases for contries, and cities, and improve existing test cases.
2. Add NOK cases for continents, countries, and cities.
3. Add OK/NOK test cases 

# Integration testing

1. Robot/Behave.
2. Setup database.
3. Create tables.
4. Use API.
5. Delete tables.
6. Remove database.
7. Deploy integration tests into the same Docker network.

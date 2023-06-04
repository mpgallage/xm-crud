# XM CRUD

XM CRUD is a Golang-based microservice that provides CRUD operations for Company objects. It uses PostgreSQL as its database and Kafka for event streaming. The application and its dependencies are containerized using Docker for seamless development and deployment.

## Features

- Company object creation, retrieval, update, and deletion.
- JWT-based authentication.
- Event streaming via Kafka.
- PostgreSQL database.

## Prerequisites

Before you begin, ensure you have Docker installed on your machine.

## Building the application

To build the application, navigate to the project directory and run the following command:

```bash
docker-compose build
```

## Running the application

To start the application, use the following command:

```bash
docker-compose up
```

## Using the services

First, obtain a JWT token by sending a request to [localhost:8080/authenticate]() with the following JSON:

```json
{
    "Username": "username",
    "Password": "password"
}
```
Note: You can replace `username` and `password` with your desired values.

You can then use this JWT token as the Bearer token in subsequent requests.

## Endpoints

Here are the available endpoints:

- `GET` [localhost:8080/company/{id}](): Retrieves a company with the given ID.
- `POST` [localhost:8080/company](): Creates a new company. The request body should look like this:

    ```json
    {
        "ID": "034c9576-ffef-48db-ab34-36fdd1fd1d45",
        "Name": "DalaLabs77",
        "Description": "This is a description of MyCompany.",
        "EmployeeCount": 150,
        "Registered": true,
        "Type": "NonProfit"
    }
    ```

- `PATCH` [localhost:8080/company/{id}](): Updates an existing company. The request body should look like this:

    ```json
    {
        "Name": "DalaLabs77",
        "Description": "This is a description of MyCompany.",
        "EmployeeCount": 150,
        "Registered": true,
        "Type": "NonProfit"
    }
    ```
  
- `DELETE` [localhost:8080/company/{id}](): Deletes an existing company.

The supported company types are: `Corporations`, `NonProfit`, `Cooperative`, and `SoleProprietorship`.

## Testing

To run integration tests, first build and run the test containers:

```bash
docker-compose -f docker-compose-test.yml build
docker-compose -f docker-compose-test.yml up
```

Then, run the tests:

```bash
go test -v ./...
```

Alternatively, you can use the `integration-tests` script to run all these steps.

## Limitations

Please note the following limitations in the current implementation:

- Username and Password are not validated against any user records. For now, any credentials provided will generate a token.

## License

XM CRUD is open source software licensed as [MIT](https://opensource.org/licenses/MIT).
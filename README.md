# CRM Backend Project
- This project supports CRUD operations to create, read, update and delete customers.
- Can also query for customer information by providing the ID as a query param to the /customers endpoint.
- All customers are structs and are currently stored in-memory in a slice of customers.
- There is static content served at /, describing the project.

## To Run
- Run either `go run main.go` or `go run *.go` from the root project directory

## Endpoints
- GET /customers
- GET /customers/{id}
- POST /customers
- PUT /customers/{id}
- DELETE /customers/{id}

### Example Request
`curl -X GET http://localhost:8000/customers`

## Future Iterations
- In the future, there are a couple of things that would be cool to do.
  - Deploy this somewhere, rather than just running on local machine.
  - Set up a database schema and write the customers to the db.
  - Docker-ize the API.
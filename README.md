## User Service
User service is a DapperLabs codding assignment.

## How to run
Provided commands will run all the environment,
apply migration and run application itself on the default port 8080

### Prerequirements
You should have docker and docker-compose installed on the machine

### Run with make file
`make rockenroll`

### Run with docker-compose
`docker-compose up -d`

## Run the tests
This command will run api tests in docker.  
Running application on port 8080 is required.

`make e2e`

## Important comments
1. I've decided to use as minimum external packages as possible in order to make the solution more explicit.
2. In order to optimize the time, only API tests on Python were added to the project.
3. I think it's better to use `PATCH` method for the partial update rather than `PUT`
   (RFC-2616 clearly mention that PUT method requests for the attached entity (in the request body) to be stored into the server).


## Authentication
JWT authentication scheme used.
Token issued during signup/login and should be hold by the client.

Authentication middleware `middleware.Authentication` holds the authentication logic.
It checks that token passed in `x-authentication-token` header,
parses it, checks expiration date and then compares user version in the token payload with
the current user version in the DB.

User version increments with every user update.

Endpoint can be excluded from the authentication flow by adding its URL
to the `noAuthUrls`

## API Specs

### `POST /signup`
Create user endpoint

**Request body**
```json
{
  "email": "test@axiomzen.co",
  "password": "123QWE!",
  "firstName": "Alex",
  "lastName": "Zimmerman"
}
```

**Constraints**

* email - should unique and valid email
* fitstName - can not be empty
* lastName - can not be empty
* password - 
  * min lenght = 5
  * should contain at least 1 digit
  * should contain at least 1 Upper case letter
  * should contain at least 1 special character

**Response**

```json
{
  "token": "some_jwt_token", "userId": 123 
}
```

**cURL**

```shell
curl -d '{"lastName":"John", "firstName":"Doe", "email": "john@doe.com", "password": "abc123A#"}' \
     -H "Content-Type: application/json" \
     -X POST http://localhost:8080/signup
```

### `POST /login`
Endpoint to log user in.

**Request body**
```json
{
  "email": "test@axiomzen.co",
  "password": "123QWE!"
}
```

**Response**

```json
{
  "token": "some_jwt_token", "userId": 123
}
```

**cURL**

```shell
curl -d '{"email": "john@doe.com", "password": "abc123A#"}' \
     -H "Content-Type: application/json" \
     -X POST http://localhost:8080/login
```

### `GET /users`
Endpoint to retrieve a json of all users. 
This endpoint requires a valid `x-authentication-token` header to be passed in with the request.

**Response**
```json
{
  "users": [
    {
      "email": "test@axiomzen.co",
      "firstName": "Alex",
      "lastName": "Zimmerman"
    }
  ]
}
```

**cURL**

```shell
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsImV4cCI6MTY3MDg3Njk3NSwianRpIjoiRmJjWG9FRmZSc1d4UExEbkpPYkNzTlZsZ1RlTWFQRVoiLCJpYXQiOjE2NjMxMDA5NzV9.6iN4b8IZlD-i_X9MxFksH0gfpkOv3n6AHRkMLvda6lI"

curl -H "Content-Type: application/json" \
     -H "x-authentication-token: ${TOKEN}" \
     -X GET http://localhost:8080/users
```

### `PUT /users/{id}`
Endpoint to update the current user `firstName` or `lastName` only. 
This endpoint requires a valid `x-authentication-token` header to be passed in.
It updates the user of the JWT being passed in.

**Request body**
```json
{
  "firstName": "NewFirstName",
  "lastName": "NewLastName"
}
```

**Response**
```json
{
  "id": 123,
  "email": "test@axiomzen.co",
  "firstName": "Alex",
  "lastName": "Zimmerman"
}
```

**cURL**

```shell
curl -d '{"firstName": "Jack", "lastName": "Doo"}' \
      -H "Content-Type: application/json" \
      -H "x-authentication-token: ${TOKEN}" \
      -X PUT http://localhost:8080/users/1
```

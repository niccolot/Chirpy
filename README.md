# Chirpy

Implementation of a twitter-like RESTAPI. The endpoints permits to register a user, login, change credentials, post read and delete a 'chirp', aka tweet. 

In the app there are implemented endopints in order to register and user, change its credentials post/delete and list chirps. 

A fittitious webhook is implemented that changes the status of the account to the premium version **chirpy red**. 

### Authorizations and authentications

A system of authentication via JWTs is implemented (with an expire time of max 24 hours), with the relative refresh tokens (with an expire time of 60 days). The refresh token is rotated every time the user credentials are changed or the JWT are refreshed and they can be revoked with the apposite endpoint. The secret key for encripting the JWTs is stored in a `.env` file in the root of the repo.

The webhook endpoints authorization is implemented via an API key that is stored in a `.env` file in the root of the repo and the users passwords are stored in the database as hashed strings.

### Database 

The database is implemented using [postgres](https://www.postgresql.org/) with Go code generated by [sqlc](https://sqlc.dev/) and the database migrations handled by [goose](https://github.com/pressly/goose).

### Errors

If an HTTP fails a response is given in the form

```json
{
    "error": "<error message with information relative to what happened and where>",
    "status code": 1 # whatever code may be 
}
```

Internal errors, e.g. fails to parse json, read some file, use some functions etc. are reported with the status code `500` and the relative message while more specific errors are listed in the endpoints below.

### Resources

```go
type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	IsChirpyred bool `json:"is_chirpy_red"`
}

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}
```

## How to use

```shell
git clone https://github.com/niccolot/Chirpy
cd Chirpy

# for hashing functions
go getgolang.org/x/crypto

# for env variables functions
go getgithub.com/joho/godotenv

# for jwt functions
go get github.com/golang-jwt/jwt/v5

# postgres go driver
go get github.com/lib/pq

# for uuids
go get github.com/google/uuid

# for database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# for installing postgres, choose whatever password you want
sudo apt install postgresql postgresql-contrib
sudo passwd postgres
```

The dependencies can be installed by running the apposite script

```shell
cd scripts
chmod +x install_dependencies.sh
./install_dependencies.sh
cd ..
```

In order to function the server has to read some environment variables, the `POLKA_API_KEY` which mimics the api key usually provided by third party companies and the `JWT_SECRET` that is used for the enription of the JWTs. This is done via the `create_env_vars.sh` script

```shell
cd scripts
chmod +x create_env_vars.sh
./create_env_vars.sh
cd ..
```
After having created the `.env` file one can load the environment variables with

```shell
source .env
```
and start the server with the provided script that before compiling the project checks if postgres is active 

```shell
cd scripts
chmod +x start.sh
./start.sh
cd ..
```

Once the server is started one can play with it sending HTTP request from a separate terminal watching the responses. One can use the linux `curl` command or the VScode extention [thunder client](https://www.thunderclient.com/) that gives a GUI for cheking the functioning of servers.

## API endpoints

* `POST /api/users`

    Allows to register a new user

    #### Request

    ```json
    {
        "password": "1234",
        "email": "walt@white.com"
    }
    ```

    #### Response

    ```json
    {
        "id": "4b15da34-2729-444e-bff6-dc95d9c7a101",
        "created_at": "2024-10-03T07:40:53.137648Z",
        "updated_at": "2024-10-03T07:40:53.137648Z",
        "email": "walt@white.com",
        "is_chirpy_red": false # no premium subscription by default
    }
    ```

    Status code: `201`

* `PUT /api/users`

    Allows to change user credentials. After this the refresh token is rotated for safety.

    #### Request

    The header must contain the users JWT

    ```
    Authorization: "Bearer <jwt>"
    ```

    ```json
    {
        "password": "newpassword",
        "email": "new@email.com"
    }
    ```

    #### Response

    ```json
    {
        "id": "4b15da34-2729-444e-bff6-dc95d9c7a101",
        "created_at": "2024-10-03T07:40:53.137648Z",
        "updated_at": "2024-10-04T07:40:53.137648Z",
        "email": "new@email.com",
        "is_chirpy_red": false
    }
    ```

    Status code: `200`

    #### Possible errors

    If the JWT is invalid the request is denied

    * Message: `invalid token`
    * Status code: `401`

* `DELETE /api/users/{id}`

    Allows to delete the user correspoinding to `{id}`. This endpoint will also delete every chirp associated with that user.

    #### Request

    The header must contain the users JWT

    ```
    Authorization: "Bearer <jwt>"
    ```
    #### Response

    Status code: `204`

    #### Possible errors

    If the `id` field of the user in case is different from the `id` associated with the JWT the request is denied 

    * Message: `invalid user`
    * Status code: `403`

    If the users JWT is invalid the request is denied

    * Message: `invalid token`
    * Status code: `403`

* `POST /api/login`

    Allows a user to login. The predefined expiration time for the JWTs is 1 hour

    #### Request

    ```json
    {
        "password": "1234",
        "email": "walt@white.com",
    }
    ```

    #### Response

    The response will contain the user informations

    ```json
    {
        "id": "4b15da34-2729-444e-bff6-dc95d9c7a101",
        "created_at": "2024-10-03T07:40:53.137648Z",
        "updated_at": "2024-10-04T07:40:53.137648Z",
        "email": "walt@white.com",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
        "refresh_token": "7b6c3b2e16b36e56d24f5d4f5c0229c72ce8379a077c010c2d8c796362f6610c",
        "is_chirpy_red": false
    }
    ```

    #### Possible errors

    If the password in the request do not correspond to the hashed password in the database the request is denied

    * Message: `unathorized access`
    * Status code: `401`

    If the email is not in the databse the request is denied

    * Message: `user <email> not found`
    * Status code: `404`

* `POST /api/chirps`

    Allows to post a chirp

    #### Request

    The header has to contain the users JWT

    ```
    Authorization: "Bearer <jwt>"
    ```

    ```json
    {
        "body": "chirp text goes here"
    }
    ```

    #### Response

     ```json
    {
        "id": "4b15da34-2729-444e-bff6-dc95d9c7a101",
        "created_at": "2024-10-03T07:40:53.137648Z",
        "updated_at": "2024-10-03T07:40:53.137648Z",
        "body": "chirp text goes here",
        "author_id": "6520a0cd-6061-41ce-a38f-ba5631758fc7" 
    }
    ```

    #### Possible errors

    If the JWT is invalid the request is denied

    * Message: `invalid token`
    * Status code: `403`

* `GET /api/chirps`

    Allows to list every chirp in the database by returning an array. It is possible to sort the chirps in ascending (default) or descending order (according to the `chirp_id`) and retrieve chirps belonging only to a certain user by using queries in the URL.

    #### Request

    Examples of valid URLs
    
    ```
    GET http://localhost:8080/api/chirps?sort=asc&author_id=2
    GET http://localhost:8080/api/chirps?sort=asc
    GET http://localhost:8080/api/chirps?sort=desc
    GET http://localhost:8080/api/chirps
    ```

    #### Response

    Using `GET http://localhost:8080/api/chirps?sort=asc&author_id=4e3936cb-09ae-44e2-a98c-303322c4f2f3`

    ```json
    [
        {
            "id":"6520a0cd-6061-41ce-a38f-ba5631758fc7",
            "body":"Gale!",
            "author_id":"4e3936cb-09ae-44e2-a98c-303322c4f2f3"
        },
        {
            "id":"d97e9629-e5c8-46c1-a099-0d6c2615f248",
            "body":"Cmon Pinkman",
            "author_id":"4e3936cb-09ae-44e2-a98c-303322c4f2f3"
        },
        {
            "id":"1a1332be-91e5-4c58-8b2e-89f08f212e98",
            "body":"Darn that fly, I just wanna cook",
            "author_id":"4e3936cb-09ae-44e2-a98c-303322c4f2f3"
        }
    ]
    ```

* `GET /api/chirps/{id}`

    Retrieves only the chirp with id `id`

    #### Request

    `GET http://localhost:8080/api/chirps/4b15da34-2729-444e-bff6-dc95d9c7a101` 

    #### Response 

    ```json
    {
        "id":"4b15da34-2729-444e-bff6-dc95d9c7a101",
        "body":"I'm the one who knocks!",
        "author_id":"1658ddd3-2a4d-4b2a-a8a9-48f52b09cc62"
    }
    ```

* `PUT /api/chirps/{id}`

    Allows to modify a chirp.

    #### Request

    ```
    Authorization: "Bearer <jwt>"
    ```

    ```json
    {
        "id": "4b15da34-2729-444e-bff6-dc95d9c7a101",
        "Body": "new text" 
    }
    ```

    #### Response

    ```json
    {
        "id": "4b15da34-2729-444e-bff6-dc95d9c7a101",
        "created_at": "2024-10-03T07:40:53.137648Z",
        "updated_at": "2024-10-04T07:40:53.137648Z",
        "body": "new text"
    }
    ```

    Status code: `200`

    #### Possible errors

    If the JWT is invalid the request is denied

    * Message: `invalid token`
    * Status code: `401`

* `POST /api/refresh`

    Allows to refresh the jwt. After the jwt has been changed the refresh token is rotated for safety.

    #### Request 

    The header of the request has to contain the refresh token (that will be changed at the end)

    ```
    Authorization: "Bearer <refresh_token>"
    ```

    #### Response 

    ```json
    {
        "token": "new.issued.jwt",
        "refresh_token": "new_rotated_refresh_token"
    }
    ```

    #### Possible errors

    If the refresh token is not found in the databse the request is denied

    * Message: `refresh token does not exists`
    * Status code: `401`

* `POST /api/revoke`

    Allows to revoke the refresh token. After this request both the refresh token and the expire time are substituded by an empty string `""`

    #### Request

    The header of the request has to contain the refresh token

    ```
    Authorization: "Bearer <refresh_token>"
    ```

    #### Response

    Status code: `204`

    #### Possible errors

    If the refresh token is not found in the databse the request is denied

    * Message: `refresh token does not exists`
    * Status code: `401`

* `DELETE /api/chirps/{chirpId}`

    Allows to delete the chirp corresponding to `chirpId`

    #### Request

    The header must contain the users JWT

    ```
    Authorization: "Bearer <jwt>"
    ```
    #### Response

    Status code: `204`

    #### Possible errors

    If the `author_id` field of the chirp in case is different from the `user_id` associated with the JWT the request is denied 

    * Message: `invalid user`
    * Status code: `403`

    If the users JWT is invalid the request is denied

    * Message: `invalid token`
    * Status code: `403`

* `POST /api/polka/webhooks`

    Allows to post a (fictitious) payment request via a webhook in order to update the account to chirpy red

    #### Request

    The header of the request has to contain the Polka API key 

    ```
    Authorization: "ApiKey <apikey>"
    ```

    ```json
    {
        # a different string will fail to update the account
        "event": "user.upgraded", 
        "data": {
            "user_id": "1a1332be-91e5-4c58-8b2e-89f08f212e98"
        }
    }
    ```

    #### Response

    Status code: `204`

    #### Possible errors

    If the API key is invalid the request is denied

    * Message: `invalid api key`
    * Status code: `401`

    If the user is not in the databse the request is denied

    * Message: `user_id <user_id> not found`
    * Status code: `404`

* `GET /api/healthz`

    Allows to check if the server is online

    #### Response

    Status code: `200`
    Body: `OK`

* `GET /admin/metrics/`

    Allows to see the number of visits to the site by rendering the HTML page `index_admin.html`

* `/admin/reset`

    Deletes all entries in the database

* `app/*`

    Renders the `index.html` file
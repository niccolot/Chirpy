# Chirpy

Implementation of a twitter like RESTAPI. The endpoints permits to register a user, login, change credentials, post read and delete a 'chirp', aka tweet.

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
        "id": 1,
        "email": "walt@white.com",
        "is_chirpy_red": false # no premium subscription by default
    }
    ```

    Status code: `201`

* `PUT /api/users`

    Allows to change user credentials. After this the refresh token is rotated for safety.

    #### Request

    The header must contain the user's JWT

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
        "password": "newpassword",
        "email": "new@email.com"
    }
    ```

    Status code: `200`

    #### Possible errors

    If the JWT is invalid the request is denied

    * Message: `invalid token`
    * Status code: `401`

* `POST /api/login`

    Allows a user to login. The request accepts an optional argument `expires_in_seconds` that sets the lifetime of the JWT. If no time or a duration greater that 1 hour is given the expire time is set to 1 hour by default.

    #### Request

    ```json
    {
        "password": "1234",
        "email": "walt@white.com",
        "expires_in_seconds":  600 # optional
    }
    ```

    #### Response

    The response will contain the user informations

    ```json
    {
        "id": 1,
        "email": "walt@white.com",
        "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
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

    The header has to contain the user's JWT

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
        "id": 1,
        "body": "chirp text goes here",
        "author_id": 2 # user_id of the author 
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

    Using `GET http://localhost:8080/api/chirps?sort=asc&author_id=1`

    ```json
    [
        {
            "id":1,
            "body":"I'm the one who knocks!",
            "author_id":1
        },
        {
            "id":2,
            "body":"Gale!",
            "author_id":1
        },
        {
            "id":3,
            "body":"Cmon Pinkman",
            "author_id":1
        },
        {
            "id":4,
            "body":"Darn that fly, I just wanna cook",
            "author_id":1
        }
    ]
    ```

* `GET /api/chirps/{id}`

    Retrieves only the chirp with id `id`

    #### Request

    `GET http://localhost:8080/api/chirps/1` 

    #### Response 

    ```json
    {
        "id":1,
        "body":"I'm the one who knocks!",
        "author_id":1
    }
    ```

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

    The header must contain the user's JWT

    ```
    Authorization: "Bearer <jwt>"
    ```
    #### Response

    Status code: `204`

    #### Possible errors

    If the `author_id` field of the chirp in case is different from the `user_id` associated with the JWT the request is denied 

    * Message: `invalid user`
    * Status code: `403`

    If the user's JWT is invalid the request is denied

    * Message: `invalid token`
    * Status code: `403`

* `POST /api/polka/webhooks`

    Allows to post a (fictitious) payment request via a webhook in order to update the account to chirpy red

    #### Request

    The header of the request has to contain the Polka API key (see `how to use` section to see how to generate one in order to mimic the keys usually given by thir part companies).

    ```
    Authorization: "ApiKey <apikey>"
    ```

    ```json
    {
        # a different string will fail to update the account
        "event": "user.upgraded", 
        "data": {
            "user_id": 3
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

* `GET /admin/metrics/*`

    Allows to see the number of visits to the site by rendering the HTML page `index_admin.html`

* `/api/reset`

    Allows to set to 0 the number of visits to the site

* `app/*`

    Renders the `index.html` file
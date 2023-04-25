# Game Library

Game Library is an application where users can find information about games they are interested in. Users will be able to see a list of games, apply filters and search game by name, navigate to game page with detailed information. If users are authenticated they can rate games and leave their comments in the discussion section. Admins and managers will be able to add, edit and delete games and all relevant information from the system.

## Main functionality:
- Registration and authorization
- Main page with the list of games
- Filtering by several game categories and search by name
- Game page with detailed information
- Rating game
- Leaving a comment
- Like a comment of other users
- Add/edit/delete games

## System Roles:
- user
- manager (can add/delete/edit games and all relevant information)
- admin (manages user roles)

## Run
```
go run . 
```

## Usage
The HTTP server runs on localhost:8081

## Endpoints 
The majority of methods expect a body with JSON value. Exclusion is createGame. The endpoint uses form data. 

### Auth
By default the user role is *user*. 

#### POST `/auth/signup`

requires *username*, *email*, *password (min=3,max=20)*

##### example req

`POST http://localhost:8081/auth/signup`

Body
```json
{
    "email":"alona.stupak@gmail.com",
    "username": "stupachek",
    "password": "qwerty"
}
```
##### res

Body
```json
{
    "message": "Sign up was successful"
}
```
#### POST `/auth/signin`
If the sign-in is successful, the response will contain JWT that the client can use to authenticate future requests. 
requires *username*, *email*, *password*

##### example req

`POST http://localhost:8081/auth/signin`

Body
```json
{
    "email":"alona.stupak@gmail.com",
    "username": "stupachek",
    "password": "qwerty"
}
```

##### res
Body
```json
{
    "message": "Sign up was successful",
    "token": "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpZCI6ImQzNmVjM2ZkLWE2M2YtNDQ3Zi1iNDhiLTIxYWRlOTgxMWI1YyIsImV4cCI6MTY4MjQ0NzU0OX0.l7YKCxPw_jwlcKvuw8WVck4B-RnM1llBRmbaJWyX304913hYXquDtro6FdNN0eRPJByMvETUJktcObKcyhyfBQ"
}
```

Other endpoints require the Authorization header with the JWT inside.

### Users 

#### GET `/users/me` 

##### example req

`GET http://localhost:8081/users/me`

##### res
Body
```json
{
    "data": {
        "id": "d36ec3fd-a63f-447f-b48b-21ade9811b5c",
        "email": "alona.stupak@gmail.com",
        "usarname": "stupachek",
        "badge_color": "",
        "role": "user"
    }
}
```

#### GET `/users/{id}` 

##### example req

`GET http://localhost:8081/users/5022e109-7f8a-419f-af0c-14109656c4d1`

##### res
Body
```json
{
    "data": {
        "id": "5022e109-7f8a-419f-af0c-14109656c4d1",
        "email": "ira.nknchn@gmail.com",
        "usarname": "IryNknchn",
        "badge_color": "",
        "role": "user"
    }
}
```

#### GET `/users` 

##### example req

`GET http://localhost:8081/users`

##### res
Body
```json
{
    "data": [
        {
            "id": "33b059c6-1f9b-43cb-9943-f04b270a6495",
            "email": "admin@a.a",
            "usarname": "admin",
            "badge_color": "",
            "role": "admin"
        },
        {
            "id": "d36ec3fd-a63f-447f-b48b-21ade9811b5c",
            "email": "alona.stupak@gmail.com",
            "usarname": "stupachek",
            "badge_color": "",
            "role": "user"
        },
        {
            "id": "5022e109-7f8a-419f-af0c-14109656c4d1",
            "email": "ira.nknchn@gmail.com",
            "usarname": "IryNknchn",
            "badge_color": "",
            "role": "user"
        }
    ]
}
```

#### PATCH `/users/{id}`
Change users role
requires *role*: *user*|*manager*|*admin*

##### example req

`PATCH http://localhost:8081/users/d36ec3fd-a63f-447f-b48b-21ade9811b5c`

Body
```json
{
   "role": "manager"
}
```
##### res

Body
```json
{
    "data": {
        "id": "d36ec3fd-a63f-447f-b48b-21ade9811b5c",
        "email": "alona.stupak@gmail.com",
        "usarname": "stupachek",
        "badge_color": "",
        "role": "manager"
    },
    "message": "User is successfully updated"
}
```

#### DELETE `/users/{id}`

##### example req

`DELETE http://localhost:8081/users/d36ec3fd-a63f-447f-b48b-21ade9811b5c`


##### res

Body
```json
{
    "message": "User is successfully deleted"
}
```
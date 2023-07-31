# Game Library

Game Library is an application where users can find information about games they are interested in. Users will be able to see a list of games, apply filters and search game by name, navigate to game page with detailed information. If users are authenticated they can rate games and leave their comments in the discussion section. Admins and managers will be able to add, edit and delete games and all relevant information from the system.

Backend https://game-library-docker.onrender.com/games

Frontend https://game-library-frontend-git-main-stupachek.vercel.app/


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

## ER Diagram
![alt text](readme-imgs/er-diagram.jpeg)

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
ROLES: ‘admin’

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

ROLES:  ‘admin’

##### example req

`DELETE http://localhost:8081/users/d36ec3fd-a63f-447f-b48b-21ade9811b5c`


##### res

Body
```json
{
    "message": "User is successfully deleted"
}
```

### Publisher

#### POST `/publishers`
ROLES: ‘manager’ | ‘admin’

##### example req

`POST http://localhost:8081/publishers`

Body
```json
{
    "name": "Nintendo"
}
```
##### res

Body
```json
{
    "data": {
        "id": "00b52e95-ef4b-4bed-9cc9-80882d30c411",
        "name": "Nintendo"
    },
    "message": "Publisher is successfully created"
}
```

#### PATCH `/publishers/{id}`

ROLES: ‘manager’ | ‘admin’

##### example req

`PATCH http://localhost:8081/publishers/00b52e95-ef4b-4bed-9cc9-80882d30c411`

Body
```json
{
    "name": "Nintendo2.0"
}
```
##### res

Body
```json
{
    "data": {
        "id": "00b52e95-ef4b-4bed-9cc9-80882d30c411",
        "name": "Nintendo2.0"
    },
    "message": "Publisher is successfully updated"
}
```
#### GET `/publishers`

##### example req

`GET http://localhost:8081/publishers`

##### res

Body
```json
{
    "data": [
        {
            "id": "00b52e95-ef4b-4bed-9cc9-80882d30c411",
            "name": "Nintendo2.0"
        },
        {
            "id": "5955136b-8f68-4474-a925-bb18d9d3627b",
            "name": "Ubisoft"
        },
        {
            "id": "f32760fb-163c-4e54-be0d-905f7b58678a",
            "name": "Sony Interactive Entertainment"
        }
    ]
}
```

#### GET `/publishers/{id}`

##### example req

`GET http://localhost:8081/publishers/00b52e95-ef4b-4bed-9cc9-80882d30c411`

##### res

Body
```json
{
    "data": {
        "id": "00b52e95-ef4b-4bed-9cc9-80882d30c411",
        "name": "Nintendo2.0"
    }
}
```

#### DELETE `/publishers/{id}`

ROLES: ‘manager’ | ‘admin’

##### example req

`DELETE http://localhost:8081/publishers/00b52e95-ef4b-4bed-9cc9-80882d30c411`

##### res

Body
```json
{
    "message": "Publisher is successfully deleted"
}
```

### Platforms

#### POST `/platforms`

ROLES: ‘manager’ | ‘admin’

##### example req

`POST http://localhost:8081/platforms`

Body
```json
{
    "name": "Playstation"
}
```
##### res

Body
```json
{
    "data": {
        "ID": "9fe14af2-2091-4eac-b072-97d468ceabc6",
        "Name": "Playstation"
    },
    "message": "Platform is successfully created"
}
```

#### GET `/platforms`

##### example req

`GET http://localhost:8081/platforms`

##### res

Body
```json
{
    "data": [
        {
            "ID": "9fe14af2-2091-4eac-b072-97d468ceabc6",
            "Name": "Playstation"
        },
        {
            "ID": "69ee7084-d6f6-44aa-a1d1-fb2e78d4d143",
            "Name": "Xbox"
        },
        {
            "ID": "e057c98d-55d5-4770-a847-5681c651c582",
            "Name": "PC"
        }
    ]
}
```

### Genres

#### POST `/genres`

ROLES: ‘manager’ | ‘admin’

##### example req

`POST http://localhost:8081/genres`

Body
```json
{
    "name": "Role-Playing"
}
```
##### res

Body
```json
{
    "data": {
        "ID": "c4c60b23-f3f7-4255-a4e2-bd14f0e7792d",
        "Name": "Role-Playing"
    },
    "message": "Genre is successfully created"
}
```
#### GET `/genres`

##### example req

`GET http://localhost:8081/genres`

##### res

Body
```json
{
    "data": [
        {
            "ID": "e7b588b0-01e5-487f-b7ef-8635a738ce36",
            "Name": "Simulation"
        },
        {
            "ID": "c4c60b23-f3f7-4255-a4e2-bd14f0e7792d",
            "Name": "Role-Playing"
        },
        {
            "ID": "3aae177c-5b71-4e83-a08f-0638f7ea874d",
            "Name": "Adventure"
        }
    ]
}
```

### Games

#### POST `/games`

requires *title*, *file*, *publisherId*
ROLES: ‘manager’ | ‘admin’

##### example req

`POST http://localhost:8081/games`

Body
```json
{
    "title": "Game",
    "description" : "Some description",
    "file": file,
    "publisherId": "699a5565-4c7e-4d18-be3e-ea04eb4f5e4d",
    "ageRestriction": 12,
    "releaseYear": 2023,
    "gernes[]": "Simulation",
    "publishers[]": "Playstation"
}
```
##### res

Body
```json
{
    "data": {
        "link": "library/Video-Game.png"
    },
    "message": "Game is successfully created"
}
```

#### GET `/games`

##### example req

`GET http://localhost:8081/games`

##### res

Body
```json
{
    "data": [
        {
            "ID": "2cab7014-27bd-4267-8222-3df345daba7b",
            "PublisherId": "699a5565-4c7e-4d18-be3e-ea04eb4f5e4d",
            "Title": "Game",
            "Description": "Some description",
            "ImageLink": "library/Video-Game.png",
            "AgeRestriction": 12,
            "ReleaseYear": 2023,
            "UpdatedAt": "0001-01-01T00:00:00Z"
        }
    ]
}
```

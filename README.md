# todo-app-backend

## How to run the app
```commandline
docker compose up --build
```

## Description
Simple REST backend application written Go. Listens on port `8080` \
Available endpoints:
- `GET /todos`
- `GET /todos/:id`
- `POST /todos`
- `PUT /todos/:id`
- `DELETE /todos/:id`
- `PUT /todos/:id/complete`
- `PUT /todos/:id/incomplete`
- `GET /todos?completed=true`
- `GET /todos?completed=false`


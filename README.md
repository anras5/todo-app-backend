# todo-app-backend

## How to run the app
```commandline
docker compose up --build
```

## Description
Simple backend application written Go. Listens on port `8080` \
Available  REST endpoints:
- `GET /todos`
- `GET /todos/:id`
- `POST /todos`
- `PUT /todos/:id`
- `DELETE /todos/:id`
- `PUT /todos/:id/complete`
- `PUT /todos/:id/incomplete`
- `GET /todos?completed=true`
- `GET /todos?completed=false`

Available GraphQL endpoint:
- `POST /graphql`

Available gRPC service on port `9000`. The `proto` file is located in the `internal/grpc/proto` directory.

## Response times of REST, GraphQL and gRPC

- REST:
```
POST 31.70272325s.
GET ONE 19.410775833s.
UPDATE 30.546839209s.
GET ALL 286.377333ms.
DELETE 29.659430041s.
```

- GRAPHQL:
```
createTodo 49.946371125s.
getTodo 38.834429125s.
updateTodo 56.15228325s.
getTodos 900.687959ms.
deleteTodo 50.838496459s.
```

- gRPC:
```

```

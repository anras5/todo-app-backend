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

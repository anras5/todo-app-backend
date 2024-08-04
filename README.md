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

## Response times of REST, GraphQL and gRPC for 100000 requests

| Operation   | REST              | GraphQL           | gRPC              |
|-------------|-------------------|-------------------|-------------------|
| Create      | 31.70272325s  | 49.946371125s  | 27.479715708s  |
| Get One     | 19.410775833s | 38.834429125s  | 19.171058792s  |
| Update      | 30.546839209s | 56.15228325s   | 26.90399975s   |
| Get All     | 286.377333ms  | 900.687959ms   | 206.215958ms   |
| Delete      | 29.659430041s | 50.838496459s  | 30.074497541s  |

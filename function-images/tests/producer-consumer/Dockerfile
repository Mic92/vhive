FROM amd64/golang:1.15.7-alpine3.13 as server

WORKDIR /app
COPY . .
CMD ["go", "run", "./server/server.go"]

FROM amd64/golang:1.15.7-alpine3.13 as client

WORKDIR /app
COPY . .
CMD ["go", "run", "./client/client.go"]
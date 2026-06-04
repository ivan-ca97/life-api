FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .
RUN CGO_ENABLED=0 go build -gcflags "all=-N -l" -o /bin/server-debug ./cmd/server
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -o /bin/server ./cmd/server

FROM alpine:3.21 AS server-debug

COPY --from=build /go/bin/dlv /usr/local/bin/dlv
COPY --from=build /bin/server-debug /bin/server
COPY --from=build /src/migrations /app/migrations

WORKDIR /app

FROM alpine:3.21 AS server

COPY --from=build /bin/server /bin/server
COPY --from=build /src/migrations /app/migrations

WORKDIR /app

CMD ["/bin/server"]

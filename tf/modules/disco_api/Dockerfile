FROM golang:1.25-bookworm as build

RUN apt update
RUN apt install -y build-essential
RUN apt-get install ca-certificates -y
RUN gcc --version

WORKDIR /code
COPY main.go go.mod go.sum ./
RUN --mount=type=cache,target="/root/.cache/go-build" go build -v -o main main.go

# docker uses alpine too
FROM debian:12.11-slim AS runner

WORKDIR /code
COPY --from=build /code/main main
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/code/main"]

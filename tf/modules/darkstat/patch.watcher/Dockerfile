FROM golang:1.27-alpine AS build

ENV CGO_ENABLED=0

WORKDIR /code
COPY main.go go.mod go.sum ./
RUN --mount=type=cache,target="/root/.cache/go-build" go build -v -o main main.go

# docker uses alpine too
FROM docker:27.3.1-cli AS runner

WORKDIR /code
COPY --from=build /code/main main
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# can't set user here because using docker with default installation

EXPOSE 8000

ENTRYPOINT ["/code/main"]
CMD [ "main" ]

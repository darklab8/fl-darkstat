FROM golang:1.23-bullseye AS dependencies

RUN apt update
RUN apt install -y build-essential
RUN apt-get install ca-certificates -y
RUN gcc --version

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download -x
RUN go install github.com/a-h/templ/cmd/templ@v0.2.747
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4

FROM dependencies AS protoc

RUN apt-get install unzip -y
WORKDIR /install
COPY docker/install_protoc.py .
RUN python3 install_protoc.py

COPY ./docker/go.mod ./docker/go.sum ./docker/installer.go ./
RUN go mod tidy
RUN go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
CMD task grpc:protoc:gateway
WORKDIR /code

FROM protoc AS build

RUN mkdir data
COPY main.go Taskfile.yml ./
COPY darkstat darkstat
COPY darkrelay darkrelay
COPY configs configs
COPY darkcore darkcore
COPY darkapi darkapi
COPY darkrpc darkrpc
COPY darkgrpc darkgrpc

# regen grpc+gateway code. Supposedly should be not changed :)
RUN task grpc:protoc:gateway

# building golang gazilion times faster
ENV GOCACHE=/root/.cache/go-build
RUN templ generate
RUN swag init --parseDependency

ARG BUILD_VERSION
RUN echo "${BUILD_VERSION}" > darkstat/settings/version.txt
RUN --mount=type=cache,target="/root/.cache/go-build" go build -v -o main main.go

FROM debian:11.6-slim AS runner
WORKDIR /code
RUN mkdir data
COPY --from=build /code/main main
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ARG BUILD_VERSION
ENV BUILD_VERSION="${BUILD_VERSION}"
CMD /code/main web
HEALTHCHECK CMD /code/main health

# test command
# docker run -v /home/naa/apps/freelancer_related/wine_prefix_freelancer_online/drive_c/Discovery:/discovery -it -e FREELANCER_FOLDER=/discovery -p 8000:8000 -e DARKSTAT_LOG_LEVEL=DEBUG -e UTILS_LOG_LEVEL=DEBUG -e DEV_ENV=true -p 8080:8080  test

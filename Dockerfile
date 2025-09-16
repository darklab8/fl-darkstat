FROM golang:1.25-bookworm AS dependencies

RUN apt update
RUN apt install -y build-essential
RUN apt-get install ca-certificates -y
RUN gcc --version

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download -x
RUN go install github.com/a-h/templ/cmd/templ@v0.3.857
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

FROM dependencies AS build

RUN mkdir data
COPY main.go Taskfile.yml ./
COPY darkstat darkstat
COPY darkrelay darkrelay
COPY configs configs
COPY darkcore darkcore
COPY darkapis darkapis
COPY darkmap darkmap

# building golang gazilion times faster
ENV GOCACHE=/root/.cache/go-build
RUN templ generate
RUN swag init --parseDependency

ARG BUILD_VERSION
RUN echo "${BUILD_VERSION}" > darkstat/settings/version.txt
RUN --mount=type=cache,target="/root/.cache/go-build" go build -v -o main main.go

FROM debian:12.11-slim AS runner
WORKDIR /code
RUN mkdir data
COPY --from=build /code/main main
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ARG BUILD_VERSION
ENV BUILD_VERSION="${BUILD_VERSION}"
ENTRYPOINT /code/main
EXPOSE 8000
CMD web
HEALTHCHECK --interval=28s CMD /code/main health

# test command
# docker run -v /home/naa/apps/freelancer_related/wine_prefix_freelancer_online/drive_c/Discovery:/discovery -it -e FREELANCER_FOLDER=/discovery -p 8000:8000 -e DARKSTAT_LOG_LEVEL=DEBUG -e UTILS_LOG_LEVEL=DEBUG -e DEV_ENV=true -p 8080:8080  test

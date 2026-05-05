
FROM golang:1.25-bookworm AS build
ENV CGO_ENABLED=0

WORKDIR /code
COPY main.go go.mod go.sum ./
RUN --mount=type=cache,target="/root/.cache/go-build" go build -v -o main main.go

FROM docker:27.3.1-cli AS runner

WORKDIR /install
RUN mkdir ~/.ssh
RUN chmod 7000 ~/.ssh
ENV change_vers=1
COPY --chmod=600 id_rsa.git.discoverygc.com /root/.ssh/id_rsa
RUN ssh-keyscan -t rsa git.discoverygc.com >> ~/.ssh/known_hosts
COPY --from=build /code/main /install/main

WORKDIR /code
ENTRYPOINT ["/install/main"]

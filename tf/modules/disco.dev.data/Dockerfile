
FROM golang:1.25-bookworm AS build
ENV CGO_ENABLED=0

WORKDIR /code
COPY main.go go.mod go.sum ./
RUN --mount=type=cache,target="/root/.cache/go-build" go build -v -o main main.go

FROM docker:27.3.1-cli AS runner

WORKDIR /install
COPY --from=build /code/main /install/main

RUN mkdir ~/.ssh
RUN chmod 700 ~/.ssh
ENV change_vers=1
COPY --chown=1001:1001 --chmod=600 id_rsa.git.discoverygc.com /root/.ssh/id_rsa
RUN ssh-keyscan -t rsa git.discoverygc.com >> ~/.ssh/known_hosts

WORKDIR /code
# RUN addgroup -g 1001 freelancer \
# && adduser -D -u 1001 -G freelancer -h /home/freelancer -s /bin/sh freelancer
# USER 1001:1001
# can't set user here because using docker

ENTRYPOINT ["/install/main"]

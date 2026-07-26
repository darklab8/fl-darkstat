FROM golang:1.23-bullseye

RUN mkdir ~/.ssh
RUN chmod 7000 ~/.ssh
COPY --chmod=600 id_rsa.files.freelancer.vanilla /root/.ssh/id_rsa
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

WORKDIR /install
COPY --chmod=777 entrypoint.sh entrypoint.sh

WORKDIR /code
CMD /install/entrypoint.sh

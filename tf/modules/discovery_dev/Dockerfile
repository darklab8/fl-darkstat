FROM docker:27.3.1-cli

RUN mkdir ~/.ssh
RUN chmod 7000 ~/.ssh
COPY --chmod=600 id_rsa.git.discoverygc.com /root/.ssh/id_rsa
RUN ssh-keyscan -t rsa git.discoverygc.com >> ~/.ssh/known_hosts
RUN ssh-keyscan -t rsa discoverygc.com >> ~/.ssh/known_hosts

WORKDIR /install
COPY --chmod=777 entrypoint.sh entrypoint.sh

WORKDIR /code
CMD /install/entrypoint.sh

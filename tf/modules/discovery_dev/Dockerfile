FROM docker:27.3.1-cli

RUN apk add curl python3 py3-pip
RUN pip3 install --break-system-packages requests

RUN mkdir ~/.ssh
RUN chmod 7000 ~/.ssh
COPY --chmod=600 id_rsa.git.discoverygc.com /root/.ssh/id_rsa
RUN ssh-keyscan -t rsa git.discoverygc.com >> ~/.ssh/known_hosts
RUN ssh-keyscan -t rsa discoverygc.com >> ~/.ssh/known_hosts

WORKDIR /install
COPY --chmod=777 entrypoint.py entrypoint.py

WORKDIR /code
CMD python3 /install/entrypoint.py

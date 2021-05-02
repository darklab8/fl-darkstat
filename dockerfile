from ubuntu:20.04

RUN adduser --disabled-password --gecos "" user
WORKDIR /home/user

RUN apt-get update
RUN apt-get -y install python3.8 python3-pip python3-venv
RUN python3 -m venv venv2

COPY . .
RUN rm -r venv
RUN venv2/bin/pip install -r requirements.txt
RUN venv2/bin/python manage.py migrate

EXPOSE 8000
CMD venv2/bin/python scripts.py manage -b run -v venv2

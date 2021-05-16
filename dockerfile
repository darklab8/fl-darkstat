from ubuntu:20.04

RUN apt-get update
RUN apt-get -y install python3.8 python3-pip python3-venv
RUN python3 -m venv venv

COPY . .
RUN venv/bin/pip install -r requirements.txt
RUN venv/bin/python manage.py migrate

EXPOSE 8000
CMD venv/bin/python scripts.py manage -b run -v venv

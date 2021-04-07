from ubuntu:20.10

RUN adduser --disabled-password --gecos "" username
WORKDIR /home/user

run apt-get update
RUN apt-get -y install python3.8 python3-pip python3-venv
RUN apt-get -y install git
RUN python3 -m venv venv

#RUN $PWD
COPY . .
RUN venv/bin/pip install -r requirements.txt
RUN venv/bin/python manage.py migrate
#RUN pip install -r requirements.txt
#RUN venv/bin/pip install hypercorn

EXPOSE 4646
RUN dir
#ENTRYPOINT ["./boot.sh"]
CMD venv/bin/python manage.py runserver --insecure 0.0.0.0:4646
#CMD python3 manage.py runserver --insecure 0.0.0.0:4646


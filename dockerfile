from ubuntu:20.04

RUN adduser user
WORKDIR /home/user

run apt-get update
RUN apt-get -y install python3.8 python3-pip
#RUN apt-get -y install git
#RUN python3 -m venv venv

#RUN $PWD
COPY . .
RUN pip3 install -r requirements.txt
#RUN python3 manage.py migrate
#RUN pip install -r requirements.txt
#RUN venv/bin/pip install hypercorn

EXPOSE 8000
RUN dir
#ENTRYPOINT ["./boot.sh"]
CMD python3 manage.py runserver --insecure 0.0.0.0:8000
#CMD python3 manage.py runserver --insecure 0.0.0.0:4646


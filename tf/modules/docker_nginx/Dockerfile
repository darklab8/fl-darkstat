FROM nginx:1.23.2-alpine

WORKDIR /code
COPY ./default.conf /etc/nginx/conf.d/default.conf
COPY ./run.sh /code/run.sh

CMD /code/run.sh

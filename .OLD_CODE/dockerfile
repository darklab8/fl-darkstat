FROM python:3.8-slim

ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

COPY ./requirements.txt ./
RUN pip install -r requirements.txt

COPY . .
RUN python manage.py migrate

EXPOSE 8000
RUN python manage.py collectstatic -c --noinput
CMD gunicorn core.wsgi -b 0.0.0.0:8000

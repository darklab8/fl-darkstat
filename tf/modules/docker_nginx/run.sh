set -e # enables automatically printed current executed commands
set -x # fails whole script (and returns error return code) if one of commands fails
# for self signed certs

apk update
apk add openssl

openssl req -x509 -nodes \
-days 365 \
-subj "/C=CA/ST=QC/O=Company, Inc./CN=127.0.0.1" \
-addext "subjectAltName=DNS:127.0.0.1" \
-newkey rsa:2048 \
-keyout /code/ssl.key \
-out /code/ssl.crt;

nginx -g 'daemon off;'

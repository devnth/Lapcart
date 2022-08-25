FROM  postgres:14-alpine 

COPY ./migrations/*.sql /docker-entrypoint-initdb.d/







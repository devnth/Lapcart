FROM postgres:14-alpine AS builder

COPY ./migrations/*.sql /docker-entrypoint-initdb.d/






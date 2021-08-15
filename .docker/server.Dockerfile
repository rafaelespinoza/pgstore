FROM postgres:12.6-alpine
WORKDIR /var/lib/postgresql
COPY .docker/server.sh scripts/
RUN scripts/server.sh
EXPOSE 5432

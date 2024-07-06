FROM golang:1.19.1
FROM ubuntu:latest

WORKDIR app
COPY ./effectiveMobile ./
COPY ./.env ./
EXPOSE $SERVER_PORT
EXPOSE $SERVER_PPROF_PORT
ENTRYPOINT ["./effectiveMobile"]
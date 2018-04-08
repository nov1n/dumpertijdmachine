FROM ubuntu:latest
COPY out/server /

RUN apt-get update
RUN apt-get install -y ca-certificates

EXPOSE 8000

CMD ["/server"]


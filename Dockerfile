FROM debian:stable-slim

RUN apt-get update \
 && apt-get install -y libreoffice-writer chromium \
 && rm -rf /var/lib/apt/lists/*

COPY setlist /setlist

RUN mkdir /repertoire

WORKDIR /repertoire

ENTRYPOINT [ "/setlist" ]

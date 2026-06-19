FROM docker.io/chromedp/headless-shell:latest

ARG TARGETPLATFORM

RUN apt-get update \
 && apt-get install -y libreoffice-writer \
 && rm -rf /var/lib/apt/lists/*

RUN mkdir /home/user \
 && chmod 777 /home/user

ENV HOME=/home/user \
    OS_ENV=container

RUN mkdir /repertoire

WORKDIR /repertoire

COPY ${TARGETPLATFORM}/setlist /setlist

ENTRYPOINT [ "/setlist" ]

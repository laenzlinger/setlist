FROM docker.io/chromedp/headless-shell:latest

RUN apt-get update \
 && apt-get install -y libreoffice-writer \
 && rm -rf /var/lib/apt/lists/*

COPY setlist /setlist

RUN mkdir /home/user \
 && chmod 777 /home/user

ENV HOME=/home/user \
    OS_ENV=container

RUN mkdir /repertoire

WORKDIR /repertoire

ENTRYPOINT [ "/setlist" ]

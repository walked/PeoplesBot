FROM ubuntu

RUN apt update && apt install -y ca-certificates
COPY PeoplesBot /app/PeoplesBot
COPY banlist.json /app/banlist.json
#resolves an issue with s3 certificate validation at authentication timedock

WORKDIR /app
ENTRYPOINT /app/PeoplesBot
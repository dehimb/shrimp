FROM alpine:latest

RUN mkdir /app
WORKDIR /app

ADD data /app/data
ADD clientapi /app/clientapi
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait /app/wait

RUN chmod +x /app/wait
RUN apk add --no-cache libc6-compat

CMD /app/wait && /app/clientapi

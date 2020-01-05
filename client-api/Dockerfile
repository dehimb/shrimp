FROM alpine:latest

RUN mkdir /app
WORKDIR /app
ADD data /app/data
ADD client-api /app/client-api
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait /app/wait
RUN chmod +x /app/wait

CMD /app/wait && /app/client-api

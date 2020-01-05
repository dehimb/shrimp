FROM alpine:latest

ADD port-domain-service /port-domain-service

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait /wait
RUN chmod +x /wait

CMD /wait && /port-domain-service

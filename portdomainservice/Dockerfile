FROM alpine:latest

ADD portdomainservice /portdomainservice
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait /wait

RUN chmod +x /wait
RUN apk add --no-cache libc6-compat

CMD /wait && /portdomainservice

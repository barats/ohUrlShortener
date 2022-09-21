FROM alpine:latest as builder
RUN apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && apk del tzdata

FROM scratch
COPY --from=builder /etc/localtime /etc/localtime
COPY ohUrlShortener /app/ohurlshortener
COPY docker/docker_config.ini /app/config.ini
EXPOSE 9091
ENTRYPOINT ["/app/ohurlshortener","-s","portal","-c","/app/config.ini"] 
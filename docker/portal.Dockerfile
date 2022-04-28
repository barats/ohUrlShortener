##
## Build 
##
FROM golang:1.16-alpine AS ohurlshortener_builder 
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.com.cn,direct
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ohurlshortener

##
## Deploy
##
FROM alpine:latest
RUN apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && apk del tzdata
WORKDIR /app
COPY --from=ohurlshortener_builder /app/ohurlshortener .
EXPOSE 9091
ENTRYPOINT ["/app/ohurlshortener","-s","portal","-c","config.ini"] 


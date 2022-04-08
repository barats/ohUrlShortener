##
## Build 
##
FROM golang:1.16-alpine AS ohurlshortener_builder 
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.com.cn,direct
ADD . /app
WORKDIR /app
# RUN go mod download && go build -o ohurlshortener .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ohurlshortener
RUN apk add upx && \
    mv ohurlshortener ohurlshortener_tmp && \
    upx --best -o ohurlshortener ohurlshortener_tmp

##
## Deploy
##
FROM alpine:latest
RUN apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && apk del tzdata
WORKDIR /app
COPY --from=ohurlshortener_builder /app/ohurlshortener .
EXPOSE 9092
ENTRYPOINT ["/app/ohurlshortener","-s","admin","-c","config.ini"] 


FROM golang:1.24-alpine

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
ENV TZ=Asia/Shanghai

RUN sed -i 's#https://dl-cdn.alpinelinux.org#https://mirrors.ustc.edu.cn#g; s#http://dl-cdn.alpinelinux.org#https://mirrors.ustc.edu.cn#g' /etc/apk/repositories \
    && apk update --no-cache \
    && apk add --no-cache tzdata ca-certificates git docker-cli \
    && go install github.com/air-verse/air@v1.62.0 \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.conf"]

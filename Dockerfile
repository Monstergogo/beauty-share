FROM golang:1.19 as build

COPY . /beauty-share
RUN if [ "$ENABLE_PROXY" = "true" ] ; then go env -w GOPROXY=https://goproxy.io,direct ; fi \
        && go env -w GO111MODULE=on \
        && go env -w CGO_ENABLED=0 \
        && cd /beauty-share/cmd/server \
        && go build -x -v -ldflags="-s -w" .

FROM alpine:latest as prod
WORKDIR /app
# 修改时区信息
RUN apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

COPY --from=build /beauty-share/cmd/server/server ./
COPY --from=build /beauty-share/conf ./conf
COPY --from=build /beauty-share/cache ./cache
# http service port
EXPOSE 5008
# grpc service port
EXPOSE 5018

CMD [ "./server" ]
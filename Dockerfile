FROM golang:1.10 as build

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x  /usr/local/bin/dep 

RUN mkdir -p /go/src/github.com/sunilkumarmohanty/site-monitor
WORKDIR /go/src/github.com/sunilkumarmohanty/site-monitor

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w" -o app cmd/main.go


FROM alpine:3.7
RUN apk add --update --no-cache ca-certificates

RUN addgroup -S monitor && adduser -S -G monitor monitor
RUN mkdir -p /home/monitor
RUN chown monitor /home/monitor

WORKDIR /home/monitor

COPY --from=build /go/src/github.com/sunilkumarmohanty/site-monitor/monitor.config .

COPY --from=build /go/src/github.com/sunilkumarmohanty/site-monitor/app .

USER monitor

CMD ["./app"]
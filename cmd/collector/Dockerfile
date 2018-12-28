FROM alpine:latest as alpine
RUN apk add --update --no-cache zip ca-certificates tzdata
WORKDIR /usr/share/zoneinfo
RUN zip -r -0 /zoneinfo.zip .

FROM scratch

# tls certificates
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# timezone data
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /

EXPOSE 14267
COPY collector-linux /go/bin/
CMD ["/go/bin/collector-linux"]

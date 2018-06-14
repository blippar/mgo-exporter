FROM golang:1.7-alpine AS builder

RUN apk add --no-cache make

ARG VERSION="unknown"
ENV MGOEXPORT_VERSION="$VERSION"

COPY .  /go/src/github.com/blippar/mgo-exporter
WORKDIR /go/src/github.com/blippar/mgo-exporter

RUN make VERSION="${MGOEXPORT_VERSION}" static

FROM scratch AS runtime

COPY --from=builder /go/src/github.com/blippar/mgo-exporter/bin/mgo-exporter /mgo-exporter

ENTRYPOINT ["/mgo-exporter"]
CMD        ["-h"]

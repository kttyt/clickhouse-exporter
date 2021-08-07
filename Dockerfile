# ===================
# ===== Builder =====
# ====================

FROM golang:1.14 AS builder
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV REPO=clickhouse-exporter

WORKDIR /src

ARG VERSION
ARG GIT_SHA
ARG NOW

ADD . .
RUN go mod tidy
RUN echo ${VERSION} ${GIT_SHA} ${NOW}
RUN go build -a \
    -ldflags " \
        -X '${REPO}/version.Version=${VERSION}' \
        -X '${REPO}/version.GitSHA=${GIT_SHA}'  \
        -X '${REPO}/version.BuiltAt=${NOW}'     \
    " \
    -o /tmp/clickhouse-exporter

# ============================
# ===== clickhouse exporter =====
# ============================

FROM alpine:3 AS clickhouse-exporter
EXPOSE 8888

RUN apk add --no-cache ca-certificates

WORKDIR /

# Copy clickhouse-exporter binary into exporter image from builder
COPY --from=builder /tmp/clickhouse-exporter .

# Run /clickhouse-exporter -alsologtostderr=true -v=1
# We can specify additional options, such as:
#   --config=/path/to/config
#   --kube-config=/path/to/kubeconf
ENTRYPOINT ["/clickhouse-exporter"]
CMD ["-logtostderr=true", "-v=1"]

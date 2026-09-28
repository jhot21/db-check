FROM --platform=$BUILDPLATFORM golang:1 AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Cross-compile natively on the build platform instead of emulating the target.
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags "-s -w -X github.com/jhot21/db-check/cmd.version=$VERSION" -o /dbcheck .

FROM busybox

COPY --from=builder /dbcheck /dbcheck

ENV TYPE=mysql

HEALTHCHECK --interval=15s --timeout=35s --start-period=30s --retries=5 CMD /dbcheck $TYPE

CMD [ "sleep", "infinity" ]

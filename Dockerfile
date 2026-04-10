FROM golang:1 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /dbcheck

FROM busybox

COPY --from=builder /dbcheck /dbcheck

ENV TYPE=mysql

HEALTHCHECK --interval=15s --timeout=35s --start-period=30s --retries=5 CMD /dbcheck $TYPE

CMD [ "sleep", "infinity" ]
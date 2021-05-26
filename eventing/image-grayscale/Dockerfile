FROM golang:1.15 as builder

COPY .  /app/
WORKDIR /app/

RUN CGO_ENABLED=0 go build -o /app/app .

FROM gcr.io/distroless/static

COPY --from=builder /app/app /app/app

CMD ["/app/app"]

FROM golang:1.15 as builder

COPY .  /app/
WORKDIR /app/

RUN CGO_ENABLED=0 go build -o /app/echo .

FROM gcr.io/distroless/static

COPY --from=builder /app/echo /app/echo

EXPOSE 8080
CMD ["/app/echo"]

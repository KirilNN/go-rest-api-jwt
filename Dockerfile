FROM golang as builder

RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/KirilNN/go-rest-api-jwt

FROM scratch

COPY --from=builder /go/bin/go-rest-api-jwt .

CMD ["./go-rest-api-jwt"]
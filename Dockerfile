FROM golang:1.19.4-alpine3.16 as Build

ARG application

WORKDIR /usr/src/app

COPY . .

RUN go build -v -o /app cmd/${application}/main.go

FROM alpine:3.16

COPY --from=Build /app /app

ARG addr
ENV addr_env=$addr

ENTRYPOINT "/app" "-addr" "${addr_env}"
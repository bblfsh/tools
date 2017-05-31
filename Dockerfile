FROM alpine:3.5

RUN apk add --no-cache device-mapper ca-certificates
ADD main bblfsh-tools
ENTRYPOINT ./bblfsh-tools
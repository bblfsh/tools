FROM alpine

RUN apk add --no-cache device-mapper ca-certificates
ADD main test
ENTRYPOINT ./test
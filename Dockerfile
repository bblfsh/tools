FROM alpine:3.5

RUN apk add --no-cache device-mapper ca-certificates
ADD build/bblfsh-tools /bin/bblfsh-tools
ADD dummy.go dummy.go
CMD /bin/bblfsh-tools dummy dummy.go
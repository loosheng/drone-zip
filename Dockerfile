FROM alpine:3.19

ADD drone-zip-release /bin/
ENTRYPOINT ["/bin/drone-zip-release"]
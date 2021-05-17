FROM alpine

ADD release/drone-zip-release /bin/
ENTRYPOINT ["/bin/drone-zip-release"]
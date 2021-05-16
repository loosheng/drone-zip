FROM plugins/base:multiarch

LABEL maintainer="LunaGod <owmurong@gmail.com>" \
  org.label-schema.name="Drone Zip Release" \
  org.label-schema.vendor="LunaGod" \
  org.label-schema.schema-version="1.0"

ADD release/drone-zip-release /bin/
ENTRYPOINT ["/bin/drone-zip-release"]
# Keycloak encourages the use of customized containers for faster startup, see
# https://www.keycloak.org/server/containers
FROM quay.io/keycloak/keycloak:23.0.3@sha256:f381819b9af906ac9f9e790ff76220c5acdeeadf50f2eeee9f861545ad6dbf1a as builder

WORKDIR /opt/keycloak
COPY cache.xml /opt/keycloak/conf/cache-ispn.xml
RUN /opt/keycloak/bin/kc.sh build \
    --db=postgres \
    --cache=ispn \
    --cache-config-file=cache-ispn.xml \
    --cache-stack=kubernetes \
    --health-enabled=true \
    --metrics-enabled=true

FROM quay.io/keycloak/keycloak:23.0.3@sha256:f381819b9af906ac9f9e790ff76220c5acdeeadf50f2eeee9f861545ad6dbf1a
COPY --from=builder /opt/keycloak/ /opt/keycloak/

# Set some reasonable default runtime options
ENV KC_LOG_CONSOLE_OUTPUT=json

ENTRYPOINT ["/opt/keycloak/bin/kc.sh"]

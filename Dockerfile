FROM caddy:2.7-builder AS builder

COPY . /tmp/src

RUN xcaddy build \
    --with github.com/patrickeasters/caddy-upstream-oauth=/tmp/src

FROM caddy:2.7

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
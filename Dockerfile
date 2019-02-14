# Build Stage
FROM lacion/alpine-golang-buildimage:1.11 AS build-stage

LABEL app="build-krontab"
LABEL REPO="https://github.com/jacobtomlinson/krontab"

ENV PROJPATH=/go/src/github.com/jacobtomlinson/krontab

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/jacobtomlinson/krontab
WORKDIR /go/src/github.com/jacobtomlinson/krontab

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/jacobtomlinson/krontab"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/krontab/bin

WORKDIR /opt/krontab/bin

COPY --from=build-stage /go/src/github.com/jacobtomlinson/krontab/bin/krontab /opt/krontab/bin/
RUN chmod +x /opt/krontab/bin/krontab

# Create appuser
RUN adduser -D -g '' krontab
USER krontab

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/krontab/bin/krontab"]

FROM openshift/origin-release:golang-1.13 AS builder
LABEL maintainer "Devtools <devtools@redhat.com>"

ENV LANG=en_US.utf8
ENV GIT_COMMITTER_NAME devtools
ENV GIT_COMMITTER_EMAIL devtools@redhat.com

ARG VERBOSE=1

WORKDIR /go/src/github.com/otaviof/chart-streams


COPY . .
RUN make build
RUN ls -ltr build/ && pwd && chmod 777 build/*

USER 10001


CMD [ "/bin/sh" ]

EXPOSE 8080
FROM openshift/origin-release:golang-1.13 AS builder
LABEL maintainer "Devtools <devtools@redhat.com>"

ENV LANG=en_US.utf8
ENV GIT_COMMITTER_NAME devtools
ENV GIT_COMMITTER_EMAIL devtools@redhat.com

ARG VERBOSE=1

WORKDIR /go/src/github.com/otaviof/chart-streams


COPY . .
RUN make build
RUN ls -ltr build/ && pwd

#--------------------------------------------------------------------

FROM registry.access.redhat.com/ubi7/ubi-minimal

LABEL com.redhat.delivery.appregistry=true
LABEL maintainer "Devtools <devtools@redhat.com>"
LABEL author "Shoubhik Bose <shbose@redhat.com>"
ENV LANG=en_US.utf8

WORKDIR /usr/local/chart-streams/bin

COPY --from=builder /go/src/github.com/otaviof/chart-streams/build/chart-streams /usr/local/chart-streams/bin/chart-streams

RUN ls -ltr /usr/local/chart-streams/bin
USER 10001

ENTRYPOINT [ "./chart-streams serve" ]

EXPOSE 8080
FROM openshift/origin-release:golang-1.13 AS builder
LABEL maintainer "Devtools <devtools@redhat.com>"
ENV GOPATH /tmp/go

ENV LANG=en_US.utf8
ENV GIT_COMMITTER_NAME devtools
ENV GIT_COMMITTER_EMAIL devtools@redhat.com

ARG VERBOSE=1

WORKDIR /go/src/github.com/otaviof/chart-streams


COPY . .
RUN make 

#--------------------------------------------------------------------

FROM registry.access.redhat.com/ubi7/ubi-minimal

LABEL com.redhat.delivery.appregistry=true
LABEL maintainer "Devtools <devtools@redhat.com>"
LABEL author "Shoubhik Bose <shbose@redhat.com>"
ENV LANG=en_US.utf8

COPY --from=builder /go/src/github.com/otaviof/chart-streams/build/chart-streams /usr/local/bin/chart-streams

RUN ls -ltr /usr/local/bin
USER 10001

WORKDIR /usr/local
CMD [ "./bin/chart-streams serve" ]

EXPOSE 8080
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

#--------------------------------------------------------------------
FROM centos:7
ENV LANG=en_US.utf8
ENV APP_INSTALL_PREFIX=/usr/local/app-server

ENV GOPATH=/tmp/go

# Create a non-root user and a group with the same name: "appserver"
ENV APP_USER_NAME=appserver
RUN useradd --no-create-home -s /bin/bash ${APP_USER_NAME}

COPY --from=builder /go/src/github.com/otaviof/chart-streams ${APP_INSTALL_PREFIX}/bin/app-server

# From here onwards, any RUN, CMD, or ENTRYPOINT will be run under the following user
USER ${APP_USER_NAME}

WORKDIR ${APP_INSTALL_PREFIX}
ENTRYPOINT [ "./bin/app-server serve" ]

EXPOSE 8080
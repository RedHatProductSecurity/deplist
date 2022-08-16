FROM registry.fedoraproject.org/fedora:35 AS builder
RUN yum install -y golang-bin
WORKDIR /go/src/project/
COPY . /go/src/project/
RUN go build -o /bin/deplist ./cmd/deplist/ 

FROM registry.fedoraproject.org/fedora:35
RUN dnf install -y \
    golang-bin \
    yarnpkg \
    maven \
    rubygem-bundler \
    ruby-devel \
    gcc \
    gcc-c++ \ 
    npm \
    && dnf clean all
COPY --from=builder /bin/deplist /
ENTRYPOINT ["/deplist"]

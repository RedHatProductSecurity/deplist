FROM golang:1.22-alpine AS build
WORKDIR /go/src/project/
COPY . /go/src/project/
RUN go build -o /bin/deplist ./cmd/deplist/ 

FROM registry.fedoraproject.org/fedora:40
RUN dnf install -y \
    golang-bin-1.22* \
    yarnpkg \
    rubygem-bundler \
    ruby-devel \
    npm \
    && dnf clean all
COPY --from=build /bin/deplist /bin
ENTRYPOINT ["/bin/deplist"]

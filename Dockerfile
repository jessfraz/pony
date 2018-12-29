FROM golang:alpine as builder
MAINTAINER Jessica Frazelle <jess@linux.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	bash \
	ca-certificates

COPY . /go/src/github.com/jessfraz/pony

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/jessfraz/pony \
	&& make static \
	&& mv pony /usr/bin/pony \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/pony /usr/bin/pony
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "pony" ]
CMD [ "--help" ]

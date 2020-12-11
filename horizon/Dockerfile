FROM golang:1.14-stretch AS builder

RUN go get github.com/stellar/go/exp/services/captivecore
RUN go install github.com/stellar/go/exp/services/captivecore

FROM ubuntu:18.04

ENV STELLAR_CORE_BINARY_PATH /usr/bin/stellar-core

ENV DEBIAN_FRONTEND=noninteractive
# ca-certificates are required to make tls connections
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl wget gnupg apt-utils
RUN wget -qO - https://apt.stellar.org/SDF.asc | APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=true apt-key add -
RUN echo "deb https://apt.stellar.org xenial stable" >/etc/apt/sources.list.d/SDF.list
RUN apt-get update && apt-get install -y stellar-core stellar-horizon
RUN apt-get clean

COPY --from=builder /go/bin/captivecore ./

ENTRYPOINT ["stellar-horizon", "--apply-migrations"]
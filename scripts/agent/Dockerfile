FROM quay.io/libpod/alpine

#RUN apk add build-base
RUN wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin
ENV GOPATH /go
ENV PATH $PATH:/go/bin

COPY . /go-daev2
WORKDIR /go-daev2

#CMD go run -race cmd/agent/main.go -h orchestrator -p 8081
CMD go run cmd/agent/main.go -h orchestrator -p 8081

FROM golang:latest
WORKDIR /message_service
COPY ./ ./
RUN go clean --modcache
RUN go mod download
RUN go get github.com/lib/pq
RUN go build -o /build gateway/cmd/main.go
EXPOSE 10001


RUN git clone https://github.com/vishnubob/wait-for-it.git && ls

CMD ["/build"]
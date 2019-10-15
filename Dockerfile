FROM golang:alpine 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app
RUN go get -d -v ./...
RUN go build -ldflags="-s -w" -o main . 
CMD ["/app/main"]
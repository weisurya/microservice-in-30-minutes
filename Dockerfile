FROM golang:1.12-alpine

# Allow Go to retrieve the dependencies for the build step
RUN apk add --no-cache git

# Secure against running as root
RUN adduser -D -u 10000 admin
RUN mkdir /gopherconuk/ && chown admin /gopherconuk/
USER admin

WORKDIR /gopherconuk/
ADD . /gopherconuk/

RUN CGO_ENABLED=0 go build -o /gopherconuk/go .

EXPOSE 8080

CMD ["/go"]
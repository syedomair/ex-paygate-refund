###########
# Stage 1 #
###########
FROM golang:1.16-alpine as builder

# Installing dependencies



# Bootstrapping modules dependencies
RUN mkdir -p /src/ex-paygate-refund
WORKDIR /src/ex-paygate-refund
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY go.mod go.mod
COPY go.sum go.sum
RUN go get -v all

# Copying source files
COPY . /src/ex-paygate-refund

# Running tests


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o api

# Create folder for scratch container
RUN mkdir /tmp/src

###########
# Stage 2 #
###########
FROM scratch


ARG service_name=""
ARG log_level=DEBUG
ARG database_url=""
ARG port=8185
ARG signingkey=""

ENV SERVICE_NAME=$service_name
ENV LOG_LEVEL=$log_level
ENV DATABASE_URL=$database_url
ENV PORT=$port
ENV SIGNINGKEY=$signingkey

EXPOSE ${PORT}

# Setup main folder
COPY --from=builder /tmp/src /src
WORKDIR /src

# Copy CA certificates which are not present in the scratch image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy api binary from first step and its dependencies
COPY --from=builder /src/ex-paygate-refund/api api

CMD ["./api"]

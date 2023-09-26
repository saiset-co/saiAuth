# Build stage
FROM golang as BUILD

WORKDIR /src/

COPY ./ /src/

RUN go build -o auth-service-bin -buildvcs=false

FROM ubuntu

WORKDIR /srv

# Copy binary from build stage
COPY --from=BUILD /src/ /srv/

RUN chmod +x /srv/auth-service-bin

# Set command to run your binary
CMD /srv/auth-service-bin start

EXPOSE 9080

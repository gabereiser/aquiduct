FROM golang:1.20.1-alpine AS build
RUN apk add --update alpine-sdk
RUN mkdir -p /srv/aquiduct
COPY . /srv/aquiduct
WORKDIR /srv/aquiduct
RUN make docker

FROM alpine:latest
RUN mkdir -p /srv/aquiduct
COPY --from=build /srv/aquiduct/build/linux/server /srv/aquiduct/server
EXPOSE 7946
WORKDIR /srv/aquiduct
CMD ./server
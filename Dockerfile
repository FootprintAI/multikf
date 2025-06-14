FROM golang:1.24-bookworm AS build
WORKDIR /src
COPY . .

WORKDIR /src
# Use go mod vendor to download imported package before building Docker image so no need to download here
#RUN go mod download

ARG GitCommitId
ARG BuildTime

# use go tool $binary | grep $variable
# to find out actual path

RUN BUILDDIR=/out make build

FROM alpine:3.22 AS bin
COPY --from=build /out /out

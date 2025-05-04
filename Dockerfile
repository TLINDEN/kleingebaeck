FROM golang:1.23-alpine as builder

RUN apk update
RUN apk upgrade
RUN apk add --no-cache git make

RUN git --version

WORKDIR /work

COPY go.mod .
COPY . .
RUN go mod download
RUN make

FROM alpine:latest
LABEL maintainer="Thomas von Dein <git@daemon.de>"

WORKDIR /app
COPY --from=builder /work/kleingebaeck /app/kleingebaeck

ENV KLEINGEBAECK_OUTDIR /backup
ENV LANG C.UTF-8
USER 1001:1001

ENTRYPOINT ["/app/kleingebaeck"]
CMD ["-h"]

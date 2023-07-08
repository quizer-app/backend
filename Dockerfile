FROM alpine:3.18.2
COPY quizer /quizer
ENTRYPOINT ["/quizer"]
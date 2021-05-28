## NOTE: This image uses goreleaser to build image
# if building manually please run go build ./cmd/yot first and then build

# Choose alpine as a base image to make this useful for CI, as many
# CI tools expect an interactive shell inside the container
FROM alpine:latest as production

#COPY --from=builder /build/yot /usr/bin/yot
COPY yot /usr/bin/yot
RUN chmod +x /usr/bin/yot

WORKDIR /workdir

ENTRYPOINT ["/usr/bin/yot"]
CMD ["--help"]
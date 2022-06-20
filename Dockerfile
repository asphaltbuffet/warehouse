FROM scratch
COPY warehouse /
ENTRYPOINT ["/warehouse"]

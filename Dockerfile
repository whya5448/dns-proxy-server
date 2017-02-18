FROM debian:8
COPY build/docker-dns-server /app/
CMD ["bash", "-c", "/app/docker-dns-server"]
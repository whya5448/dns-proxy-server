FROM debian:8
COPY build/ /app/
WORKDIR /app
VOLUME ["/var/run/docker.sock", "/var/run/docker.sock"]
CMD ["bash", "-c", "/app/dns-proxy-server"]
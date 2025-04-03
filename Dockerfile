FROM debian:stable-slim

COPY chirpy /bin/chirpy
CMD ["/bin/chirpy"]
FROM scratch
COPY confd /

ENTRYPOINT ["/confd"]
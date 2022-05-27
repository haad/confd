FROM scratch

COPY ./bin/confd /

ENTRYPOINT ["/confd"]
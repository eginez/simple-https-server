FROM frolvlad/alpine-glibc

EXPOSE 8443

COPY main .
COPY server.crt .
COPY server.key .
CMD ["./main"]


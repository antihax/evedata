FROM scratch

ADD templates /templates
ADD static /static
ADD evedata-server /bin/

ENTRYPOINT ["/bin/evedata-server"]

EXPOSE 3000

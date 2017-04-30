FROM scratch

ADD templates /templates
ADD static /static
ADD evedata-server /

ENTRYPOINT ["/evedata-server"]

EXPOSE 3000

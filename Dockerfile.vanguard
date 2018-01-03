FROM scratch

ADD services/vanguard/templates /templates
ADD services/vanguard/static /static
ADD bin/vanguard /

ENTRYPOINT ["/vanguard"]

EXPOSE 3000

FROM alpine:3.7

ADD server /
ADD docker-entrypoint.sh /

ENTRYPOINT ["sh", "docker-entrypoint.sh"]
CMD ["run"]

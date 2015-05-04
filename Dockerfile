FROM debian
MAINTAINER Sam Rayment samrayment@gmail.com

COPY monidash.tmp /monidash
ENTRYPOINT ["/monidash"]

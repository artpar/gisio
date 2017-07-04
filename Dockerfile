FROM busybox

MAINTAINER Parth Mudgal <artpar@gmail.com>
WORKDIR /opt/gisio

ADD main /opt/gisio/gisio
RUN chmod +x /opt/gisio/gisio
ADD resources /opt/gisio/resources

VOLUME /opt/gisio/data


EXPOSE 2299

ENTRYPOINT ["/opt/gisio/gisio", "data"]
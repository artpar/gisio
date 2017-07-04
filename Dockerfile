FROM busybox

MAINTAINER Parth Mudgal <artpar@gmail.com>
WORKDIR /opt/gisio

ADD main /opt/gisio/gisio
RUN chmod +x /opt/gisio/gisio
#ADD gomsweb/dist /opt/goms/gomsweb/dist

EXPOSE 8080
ENTRYPOINT ["/opt/gisio/gisio"]
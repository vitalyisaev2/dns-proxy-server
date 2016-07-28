FROM node:5.6-wheezy

ENV WORKDIR=/opt/docker-dns-server
COPY . $WORKDIR
EXPOSE 53
WORKDIR $WORKDIR
CMD /opt/docker-dns-server/discover-ip && node app.js

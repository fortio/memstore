# Go releaser dockerfile
FROM scratch
COPY memstore /usr/bin/memstore
EXPOSE 8079
EXPOSE 8080
# configmap (dynamic flags)
VOLUME /etc/memstore
# data files etc
VOLUME /var/lib/memstore
WORKDIR /var/lib/memstore
ENTRYPOINT ["/usr/bin/memstore"]
CMD ["-config-dir", "/etc/memstore"]

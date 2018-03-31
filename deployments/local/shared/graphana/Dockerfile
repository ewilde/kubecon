FROM grafana/grafana

RUN apt-get update && apt-get -y install curl

# Change the default data directory (otherwise grafana.db won't persist)
RUN mkdir /var/lib/grafanadb
ENV GF_PATHS_DATA /var/lib/grafanadb

RUN mkdir -p /var/lib/grafana/dashboards
COPY ./dashboards.yml  /etc/grafana/provisioning/dashboards/dashboards.yml
COPY ./dashboard.home.json  /var/lib/grafana/dashboards/home.json
COPY ./dashboard.home.json  /usr/share/grafana/public/dashboards/home.json
COPY ./dashboard.status.json  /var/lib/grafana/dashboards/status.json


# Init Grafana sqlite db and preconfigure our data source to be our influxdb kubecon db
RUN bash -c '/run.sh & sleep 15 && curl -s -H "Content-Type: application/json" -X POST \
    --data '"'"'{"name": "myinfluxdb", "type": "influxdb", "access": "proxy", "url": "http://influxdb:8086", "database": "kubecon", "isDefault": true}'"'"' \
    http://admin:admin@localhost:3000/api/datasources \
    && kill -SIGINT %%'

CMD ["/run.sh"]

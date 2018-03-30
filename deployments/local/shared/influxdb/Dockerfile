FROM influxdb:1.2.2

# Create a kubecon db
RUN /bin/bash -c "influxd run & sleep 5 && influx -execute 'CREATE DATABASE kubecon' && kill %1 && sleep 5"

CMD ["influxd"]

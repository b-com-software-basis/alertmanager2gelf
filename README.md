# alertmanager2gelf - Prometheus Alertmanager webhook to gelf

**alertmanager2gelf** is an [Alertmanager Webhook Receiver]
compliant with a [GELF] server such as Graylog.

[GELF]: http://docs.graylog.org/en/3.0/pages/gelf.html
[Alertmanager Webhook Receiver]: https://prometheus.io/docs/operating/integrations/#alertmanager-webhook-receiver

## Build

Build **alertmanager2gelf** with:

    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" alertmanager2gelf.go

## Installation

Copy *alertmanager2gelf* binary to /usr/local/bin

## Configuration

Create the following file /etc/alertmanager2gelf/alertmanager2gelf.yml

    listenOn: "localhost:5001"
    graylogAddr: "192.168.0.23:12201"

* listenOn: listening address and port of the service (default: "localhost:5001")
* graylogAddr: graylog compliant remote service address and port (default: "localhost:12201")

## Startup

You can use the systemd file *alertmanager2gelf.service*.

Default user is *alertmanager*, you may want to adapt this.

Copy alertmanager2gelf.service to /etc/systemd/system/alertmanager2gelf.service

Then:

    systemctl daemon-reload
    systemctl start alertmanager2gelf
    systemctl status alertmanager2gelf


## Alertmanager integration

Add *alertmanager2gelf* receiver to alertmanager.

In /etc/alertmanager/alertmanager.yml file:

    receivers:
    - name: 'my-team'
        webhook_configs:
        - send_resolved: true
            url: 'http://localhost:5001'

Then restart alertmanager service.

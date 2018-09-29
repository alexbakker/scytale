# scyserver

```sh
# add a new user for the server
useradd --system --user-group scyserver --create-home --home-dir /var/lib/scyserver --shell /sbin/nologin

# copy the service file over after reviewing it
cp etc/scyserver.service /etc/systemd/system

# start the service
systemctl daemon-reload && systemctl enable --now scyserver
```

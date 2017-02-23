# scyserver

```bash
# add a new user for the server
sudo useradd --system --user-group scyserver --create-home --home-dir /var/lib/scyserver --shell /sbin/nologin

# copy the service file over after reviewing it
sudo cp etc/scyserver.service /etc/systemd/system

# run daemon-reload and start the service
```

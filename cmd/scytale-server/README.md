# scytale-server

```bash
# add a new user for the server
sudo useradd --system --user-group scytale-server --create-home --home-dir /var/lib/scytale-server --shell /sbin/nologin

# copy the service file over after reviewing it
sudo cp etc/scytale-server.service /etc/systemd/system

# run daemon-reload and start the service
```

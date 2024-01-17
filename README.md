# dbus-snotify

Notifications, with sound.

A fork of [Kimiblock/snotify](https://github.com/Kimiblock/snotify).

## Build

```bash
# Install go 1.21
make
```
## Deploy

### Docker

```bash
docker run \
    --restart=always \
    -d \
    --name dbus-snotify \
    --volume /run/user/$(id -u)/bus:/run/user/$(id -u)/bus \
    --userns keep-id \
    -v /run/user/$(id -u)/pulse:/run/user/$(id -u)/pulse \
    -v ${XDG_RUNTIME_DIR}/pulse/native:${XDG_RUNTIME_DIR}/pulse/native \
    -e PULSE_SERVER=unix:${XDG_RUNTIME_DIR}/pulse/native \
    -e DBUS_SESSION_BUS_ADDRESS \
    naturalselect/dbus-snotifypod:latest
```

### Podman

**Start container:**

```bash
docker run \
    --restart=always \
    -d \
    --name dbus-snotify \
    --volume /run/user/$(id -u)/bus:/run/user/$(id -u)/bus \
    --userns keep-id \
    -v /run/user/$(id -u)/pulse:/run/user/$(id -u)/pulse \
    -v ${XDG_RUNTIME_DIR}/pulse/native:${XDG_RUNTIME_DIR}/pulse/native \
    -e PULSE_SERVER=unix:${XDG_RUNTIME_DIR}/pulse/native \
    -e DBUS_SESSION_BUS_ADDRESS \
    naturalselect/dbus-snotifypod:latest
```

**Generate systemd service:**

```bash
docker generate systemd --new --name dbus-snotify -f

mv container-dbus-snotify.service ~/.config/systemd/user/
```

**Enable service:**

```bash
systemctl --user enable container-dbus-snotify
systemctl --user restart container-dbus-snotify
```
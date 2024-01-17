# Build project
FROM golang:1.21 AS build
COPY . /root/dbus-notify
WORKDIR /root/dbus-notify
RUN make

# Build image
FROM ubuntu:22.04
LABEL "Maintainer"="NaturalSelect<2145973003@qq.com>"

RUN apt update
RUN apt install sudo dbus pulseaudio -y
RUN rm -rf /var/lib/apt/lists/*

RUN chown root:root /usr/bin/sudo && chmod 4755 /usr/bin/sudo
ENV HOME /home/snotify
RUN useradd --create-home --home-dir $HOME snotify
RUN echo 'snotify ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

COPY --from=build /root/dbus-notify/build/dbus-snotify ${HOME}/dbus-snotify
COPY --from=build /root/dbus-notify/message.ogg ${HOME}/message.ogg
RUN chown -R snotify:snotify $HOME
WORKDIR "/home/snotify/"
ENTRYPOINT [ "/home/snotify/dbus-snotify" ]
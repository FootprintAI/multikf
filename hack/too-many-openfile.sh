#!/usr/bin/env bash

# when running too-many-open-file issue, use
sudo sysctl fs.inotify.max_user_instances=1280
sudo sysctl fs.inotify.max_user_watches=655360

# or modify /etc/sysctl.conf to persist a reboot
# sudo vim /etc/sysctl.conf
# + fs.inotify.max_user_instances=1280
# + fs.inotify.max_user_watches=655360
#

#!/usr/bin/env bash

# when running too-many-open-file issue, use
sudo sysctl fs.inotify.max_user_instances=1280
sudo sysctl fs.inotify.max_user_watches=655360

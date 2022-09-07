#!/usr/bin/env bash

# reference: https://stackoverflow.com/questions/39293441/needed-ports-for-kubernetes-cluster
cat << EOF | tee /etc/ufw/user.rules
*filter
:ufw-user-input - [0:0]
:ufw-user-output - [0:0]
:ufw-user-forward - [0:0]
:ufw-before-logging-input - [0:0]
:ufw-before-logging-output - [0:0]
:ufw-before-logging-forward - [0:0]
:ufw-user-logging-input - [0:0]
:ufw-user-logging-output - [0:0]
:ufw-user-logging-forward - [0:0]
:ufw-after-logging-input - [0:0]
:ufw-after-logging-output - [0:0]
:ufw-after-logging-forward - [0:0]
:ufw-logging-deny - [0:0]
:ufw-logging-allow - [0:0]
:ufw-user-limit - [0:0]
:ufw-user-limit-accept - [0:0]
### RULES ###
### tuple ### allow tcp 80 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 80 -j ACCEPT
### tuple ### allow tcp 443 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 443 -j ACCEPT
### tuple ### allow tcp 22 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 22 -j ACCEPT
### tuple ### allow tcp 179 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 179 -j ACCEPT
### tuple ### allow udp 1194 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p udp --dport 1194 -j ACCEPT
### tuple ### allow tcp 8443 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 8443 -j ACCEPT
### tuple ### allow tcp 4149 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 4149 -j ACCEPT
### tuple ### allow tcp 9099 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 9099 -j ACCEPT
### tuple ### allow tcp 6443 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp --dport 6443 -j ACCEPT
### tuple ### allow tcp 2379:2380 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp -m multiport --dports 2379:2380 -j ACCEPT
### tuple ### allow tcp 10250:20256 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp -m multiport --dports 10250:20256 -j ACCEPT
### tuple ### allow tcp 30000:32767 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p tcp -m multiport --dports 30000:32767 -j ACCEPT
### tuple ### allow udp 8285 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p udp --dport 8285 -j ACCEPT
### tuple ### allow udp 8472 0.0.0.0/0 any 0.0.0.0/0 in
-A ufw-user-input -p udp --dport 8472 -j ACCEPT
### END RULES ###
### LOGGING ###
-A ufw-after-logging-input -j LOG --log-prefix "[UFW BLOCK] " -m limit --limit 3/min --limit-burst 10
-A ufw-after-logging-forward -j LOG --log-prefix "[UFW BLOCK] " -m limit --limit 3/min --limit-burst 10
-I ufw-logging-deny -m conntrack --ctstate INVALID -j RETURN -m limit --limit 3/min --limit-burst 10
-A ufw-logging-deny -j LOG --log-prefix "[UFW BLOCK] " -m limit --limit 3/min --limit-burst 10
-A ufw-logging-allow -j LOG --log-prefix "[UFW ALLOW] " -m limit --limit 3/min --limit-burst 10
### END LOGGING ###
### RATE LIMITING ###
-A ufw-user-limit -m limit --limit 3/minute -j LOG --log-prefix "[UFW LIMIT BLOCK] "
-A ufw-user-limit -j REJECT
-A ufw-user-limit-accept -j ACCEPT
### END RATE LIMITING ###
COMMIT
EOF

ufw enable

# check status
ufw status numbered


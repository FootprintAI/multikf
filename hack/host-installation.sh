#!/usr/bin/env bash

# see https://www.vagrantup.com/downloads

# macOS
brew install vagrant

# linux

curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update && sudo apt-get install vagrant

# windows
curl -Lo vagrant.msi https://releases.hashicorp.com/vagrant/2.2.19/vagrant_2.2.19_x86_64.msi

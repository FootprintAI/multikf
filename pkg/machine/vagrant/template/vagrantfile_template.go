package template

import (
	"fmt"
	"html/template"
	"io"

	pkgtemplate "github.com/footprintai/multikf/pkg/template"
)

func NewDefaultVagrantTemplate() *DefaultVagrantFileTemplate {
	return &DefaultVagrantFileTemplate{
		vagrantFileTemplate: vagrantFileDefaultTemplate,
	}
}

func (d *DefaultVagrantFileTemplate) Filename() string {
	return "Vagrantfile"
}

func (d *DefaultVagrantFileTemplate) Execute(w io.Writer) error {
	tmpl, err := template.New("vagrantfile").Parse(d.vagrantFileTemplate)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(w, d); err != nil {
		return err
	}
	return nil
}

type vagrantConfig interface {
	pkgtemplate.NameGetter
	pkgtemplate.KubeAPIPortGetter
	pkgtemplate.SSHPortGetter
	pkgtemplate.CpuMemoryGetter
}

func (d *DefaultVagrantFileTemplate) Populate(config interface{}) error {
	if _, isVagrantConfig := config.(vagrantConfig); !isVagrantConfig {
		return fmt.Errorf("config didn't implement vagrantConfig interface")
	}
	v := config.(vagrantConfig)
	d.VMName = v.GetName()
	d.KubeAPIPort = v.GetKubeAPIPort()
	d.SSHPort = v.GetSSHPort()
	d.Memory = v.GetMemory()
	d.CPUs = v.GetCPUs()
	return nil
}

type DefaultVagrantFileTemplate struct {
	KubeAPIPort         int
	SSHPort             int
	VMName              string
	Memory              int // in bytes
	CPUs                int
	vagrantFileTemplate string
}

var vagrantFileDefaultTemplate string = `
Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/focal64"
  config.vm.provision "file", source: "kind-config.yaml", destination: "/tmp/kind-config.yaml"
  config.vm.provision :shell, path: "bootstrap/bootstrap.sh"
  config.vm.provision :shell, path: "bootstrap/provision-cluster.sh"
  config.vm.provision :shell, path: "bootstrap/provision-kf14.sh"
  config.vm.network :forwarded_port, guest: {{.KubeAPIPort}}, guest_ip: "0.0.0.0", host: {{.KubeAPIPort}}
  config.vm.network :forwarded_port, guest: 22, host: {{.SSHPort}}, id: "ssh"

  # define vm name
  config.vm.define :{{.VMName}} do |t|
  end

  config.vm.provider "virtualbox" do |vb|
    # Display the VirtualBox GUI when booting the machine
    #vb.gui = true

    # Customize the amount of memory on the VM:
    vb.memory = "{{.Memory}}"
    vb.cpus = {{.CPUs}}
  end
end
`

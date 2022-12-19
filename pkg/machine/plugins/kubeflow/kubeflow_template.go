package template

import (
	"fmt"
	"html/template"
	"io"

	pkgtemplate "github.com/footprintai/multikf/pkg/template"
	"golang.org/x/crypto/bcrypt"
)

func NewKubeflowTemplateExecutor(filename string, manifest string) pkgtemplate.TemplateExecutor {
	return newKubeflowTemplateWithTemplate(filename, manifest)
}

func newKubeflowTemplateWithTemplate(filename, tmpl string) *KubeflowFileTemplate {
	return &KubeflowFileTemplate{
		filename:                filename,
		kubeflowFileTemplate:    tmpl,
		DefaultSaltedPassword:   "",
		AuthServicePVCSizeInG:   10,
		KatibMySQLPVCSizeInG:    10,
		PipelineMinioPVCSizeInG: 20,
		PipelineMySQLPVCSizeInG: 20,
	}
}

func mustBcryptGenerated(originPasswrod string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(originPasswrod), bcrypt.DefaultCost)
	return string(hashedPassword)
}

type KubeflowFileTemplate struct {
	filename                string
	DefaultSaltedPassword   string
	AuthServicePVCSizeInG   int
	KatibMySQLPVCSizeInG    int
	PipelineMinioPVCSizeInG int
	PipelineMySQLPVCSizeInG int

	kubeflowFileTemplate string
}

func (k *KubeflowFileTemplate) Filename() string {
	return k.filename
}

func (k *KubeflowFileTemplate) Execute(w io.Writer) error {
	tmpl, err := template.New("kubeflowconfig").Delims("[[", "]]").Parse(k.kubeflowFileTemplate)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(w, k); err != nil {
		return err
	}
	return nil
}

type kubeflowConfig interface {
	pkgtemplate.DefaultPasswordGetter
}

func (k *KubeflowFileTemplate) Populate(v interface{}) error {
	if _, isConfiger := v.(kubeflowConfig); !isConfiger {
		return fmt.Errorf("not implements kubeflowConfig interface")
	}
	c := v.(kubeflowConfig)
	k.DefaultSaltedPassword = mustBcryptGenerated(c.GetDefaultPassword())
	return nil
}

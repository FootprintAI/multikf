package template

import (
	"html/template"
	"io"
)

func NewAuditPolicyTemplate() *AuditPolicyFileTemplate {
	return &AuditPolicyFileTemplate{
		auditPolicyFileTemplate: auditPolicyDefaultFileTemplate,
	}
}

func (k *AuditPolicyFileTemplate) Filename() string {
	return "audit-policy.yaml"
}

func (k *AuditPolicyFileTemplate) Execute(w io.Writer) error {
	tmpl, err := template.New("auditpolicy").Parse(k.auditPolicyFileTemplate)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(w, k); err != nil {
		return err
	}
	return nil
}

func (k *AuditPolicyFileTemplate) Populate(v interface{}) error {
	return nil
}

type AuditPolicyFileTemplate struct {
	auditPolicyFileTemplate string
}

var auditPolicyDefaultFileTemplate string = `
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
- level: Metadata
`

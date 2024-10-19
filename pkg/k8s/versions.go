package k8s

import "fmt"

type KindK8sVersion struct {
	version string
	sha256  string
}

func NewKindK8sVersion(version string, sha256 string) KindK8sVersion {
	return KindK8sVersion{version: version, sha256: sha256}

}

func (k KindK8sVersion) Version() string {
	return k.version
}

func (k KindK8sVersion) Sha256() string {
	return k.sha256
}

func (k KindK8sVersion) String() string {
	return fmt.Sprintf("kindest/node:%s@sha256:%s", k.version, k.sha256)
}

func DefaultVersion() KindK8sVersion {
	return v12813
}

var (
	v1310  = NewKindK8sVersion("v1.31.0", "53df588e04085fd41ae12de0c3fe4c72f7013bba32a20e7325357a1ac94ba865")
	v1304  = NewKindK8sVersion("v1.30.4", "976ea815844d5fa93be213437e3ff5754cd599b040946b5cca43ca45c2047114")
	v1298  = NewKindK8sVersion("v1.29.8", "d46b7aa29567e93b27f7531d258c372e829d7224b25e3fc6ffdefed12476d3aa")
	v12813 = NewKindK8sVersion("v1.28.13", "45d319897776e11167e4698f6b14938eb4d52eb381d9e3d7a9086c16c69a8110")
	v12716 = NewKindK8sVersion("v1.27.16", "3fd82731af34efe19cd54ea5c25e882985bafa2c9baefe14f8deab1737d9fabe")
	v12615 = NewKindK8sVersion("v1.26.15", "1cc15d7b1edd2126ef051e359bf864f37bbcf1568e61be4d2ed1df7a3e87b354")
)

func ListVersion() []KindK8sVersion {
	return []KindK8sVersion{
		v1310,
		v1304,
		v1298,
		v12813,
		v12716,
		v12615,
	}
}

func ListVersionString() []string {
	var vs []string
	for _, v := range ListVersion() {
		vs = append(vs, v.Version())
	}
	return vs
}

func ListVersionSha256String() []string {
	var vs []string
	for _, v := range ListVersion() {
		vs = append(vs, fmt.Sprintf("%s:sha256:%s", v.Version(), v.Sha256()))
	}
	return vs
}

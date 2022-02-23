package fs

import (
	"io/fs"
	"io/ioutil"
	"os"

	//log "github.com/golang/glog"
	"github.com/spf13/afero"
)

func NewFolder(folderPath string) *Folder {
	return &Folder{
		folderFs:   afero.NewBasePathFs(afero.NewOsFs(), folderPath),
		folderPath: folderPath,
	}
}

type Folder struct {
	folderFs   afero.Fs
	folderPath string
	//forceOverwrite bool
}

func (v *Folder) Root() string {
	return v.folderPath
}

func (v *Folder) ensureFolder() error {
	_, err := v.folderFs.Stat(".")
	if os.IsNotExist(err) {
		if err := v.folderFs.MkdirAll(".", 0777); err != nil {
			return err
		}
	}
	return nil
}

// DumpFiles dumps virtualFs into a real file folder
func (v *Folder) DumpFiles(virtualFs ...fs.FS) error {
	if err := v.ensureFolder(); err != nil {
		return err
	}

	for _, vfs := range virtualFs {
		err := fs.WalkDir(vfs, ".", func(path string, d fs.DirEntry, err error) error {
			if path == "." {
				return nil
			}
			if d.IsDir() {
				//fmt.Printf("D %s\n", path)
				//os.Mkdir(path)
				if err := v.folderFs.MkdirAll(path, 0777); err != nil {
					return err
				}
				return nil
			}
			//fmt.Printf("F %s\n", path)
			fd, err := vfs.Open(path)
			if err != nil {
				return err
			}
			b, err := ioutil.ReadAll(fd)
			if err != nil {
				return err
			}
			fd.Close()
			if err := afero.WriteFile(v.folderFs, path, b, 0666); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

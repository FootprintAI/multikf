package fs

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	//log "github.com/golang/glog"
	"github.com/footprintai/multikf/pkg/machine/fsutil"
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

func (v *Folder) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(v.folderFs, filename, data, perm)
}

func (v *Folder) Exists(filename string) bool {
	return fsutil.Exists(v.IOFS(), filename)
}

func (v *Folder) IOFS() fs.FS {
	return afero.NewIOFS(v.folderFs)
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
func (v *Folder) DumpFiles(overwrite bool, virtualFs ...fs.FS) error {
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
			if v.Exists(path) && !overwrite {
				return fmt.Errorf("dumpfile: destination(%s) exists", path)
			}
			fd, err := vfs.Open(path)
			if err != nil {
				return err
			}
			b, err := ioutil.ReadAll(fd)
			if err != nil {
				return err
			}
			fd.Close()
			// TODO: check the file exists or not
			if err := v.WriteFile(path, b, 0666); err != nil {
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

package local_volume

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/idtools"
)

func (v *localVolume) allocateVolumeOnDisk() error {
	path := v.GetPath()
	rootUID, rootGID, err := idtools.GetRootUIDGID(nil, nil)
	if err != nil {
		return err
	}
	if err := idtools.MkdirAllAs(path, 0755, rootUID, rootGID); err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("volume already exists under %s", filepath.Dir(path))
		}
		return err
	}
	return nil
}

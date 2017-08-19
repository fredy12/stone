package local_volume

import (
	"fmt"
	"os"
	"path/filepath"
	"errors"
	"io/ioutil"

	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/idtools"
)

func (v *LocalVolume) allocateVolumeOnDisk(path string) error {
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

func fromDisk(pth string) (*LocalVolume, error) {
	filename := filepath.Join(pth, "volume.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, errors.New("error restore from disk without volume.json")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var obj LocalVolume
	err = json.Unmarshal(data, &obj)
	if err != nil {
		logrus.Debugf("parse %s failure: %s", filename, err)
		return nil, err
	}
	//if v.Name != obj.Name {
	//	return fmt.Errorf("name error: %s != %s", v.Name, obj.Name)
	//}
	dataPath := filepath.Join(pth, VolumeDataPathName)
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		return nil, errors.New("error restore from disk without _data path")
	}

	if obj.Name == "" || obj.VolumePath == "" || obj.DataPath == "" || obj.DiskId == "" {
		return nil, errors.New("error restore from disk with wrong format data")
	}
	return &obj, nil
}

func (v *LocalVolume) toDisk(pth string) error {
	v.m.Lock()
	defer v.m.Unlock()

	data, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(pth, "volume.json"), data, 0600); err != nil {
		return err
	}
	return nil
}

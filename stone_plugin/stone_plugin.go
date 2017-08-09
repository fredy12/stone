// Copyright (c) 2015-2016, ZANECLOUD CORPORATION. All rights reserved.

package stone_plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/zanecloud/stone/stone_plugin/tools"
	"github.com/zanecloud/stone/stone_plugin/volume"
)

const (
	DefaultDriverName = "stone"
)

var (
	ErrVolumeBadFormat   = errors.New("bad volume format")
	ErrVolumeUnsupported = errors.New("unsupported volume")
	ErrVolumeNotFound    = errors.New("no such volume")
	ErrVolumeVersion     = errors.New("invalid volume version")
	ErrVolumePathExist   = errors.New("volume path is exist")
	// ErrNotFound is the typed error returned when the requested volume name can't be found
	ErrNotFound = fmt.Errorf("volume not found")
)

type stonePlugin struct {
	m       sync.Mutex
	scope   string
	volumes map[string]volume.Volume
}

func assert(err error) {
	if err != nil {
		logrus.Errorf("Error: %v", err)
	}
}

func New() *stonePlugin {
	// TODO: add a task to auto check and load disks
	diskInfos, err := tools.Collect()
	if err != nil {
		logrus.Panicf("error when do init stone plugin with collect disks: %v", err)
	}
	diskJson, err := json.Marshal(&diskInfos)
	if err != nil {
		logrus.Panicf("error when do init stone plugin with json disks: %v", err)
	}
	logrus.Infof("Collected Disk Information: %s", string(diskJson))

	// restore existed volumes
	vols, err := restoreVolumes(diskInfos)
	if err != nil {
		logrus.Panicf("error when do init stone plugin with restore volumes: %v", err)
	}
	volumeJson, err := json.Marshal(&vols)
	if err != nil {
		logrus.Panicf("error when do init stone plugin with json volumes: %v", err)
	}
	logrus.Infof("Restore volumes Information: %s", string(volumeJson))

	return &stonePlugin{
		volumes: vols,
	}
}

func (s *stonePlugin) implement() string { return "VolumeDriver" }

func (s *stonePlugin) register(api *PluginAPI) {
	prefix := "/" + s.implement()

	api.Handle("POST", prefix+".Create", s.create)
	api.Handle("POST", prefix+".Remove", s.remove)
	api.Handle("POST", prefix+".Mount", s.mount)
	api.Handle("POST", prefix+".Unmount", s.unmount)
	api.Handle("POST", prefix+".Path", s.path)
	api.Handle("POST", prefix+".Get", s.get)
	api.Handle("POST", prefix+".List", s.list)
}

func fmtError(err error, vol string) *string {
	s := fmt.Sprintf("%v: with volume name %s", err, vol)
	return &s
}

// Name returns the name of Root, defined in the volume package in the DefaultDriverName constant.
func (s *stonePlugin) Name() string {
	return DefaultDriverName
}

func (s *stonePlugin) create(resp http.ResponseWriter, req *http.Request) {
	var q struct {
		Name string
		Opts map[string]string
	}
	var r struct{ Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received create request for volume '%s'\n", q.Name)

	if err := s.validateName(q.Name); err != nil {
		r.Err = fmtError(err, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	s.m.Lock()
	defer s.m.Unlock()

	_, exists := s.volumes[q.Name]
	if exists {
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	opts, err := s.setOpts(q.Opts)
	if err != nil {
		r.Err = fmtError(err, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	d, err := s.chooseDisk(opts)
	if err != nil {
		r.Err = fmtError(err, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	v, err := volume.New(s.Name(), q.Name, opts.size, opts.ioClass, opts.exclusive, d)
	if err != nil {
		r.Err = fmtError(err, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	s.volumes[q.Name] = v
	assert(json.NewEncoder(resp).Encode(r))
}

func (s *stonePlugin) remove(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Err *string }
	s.m.Lock()
	defer s.m.Unlock()

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received remove request for volume '%s'\n", q.Name)

	v, exists := s.volumes[q.Name]
	if !exists {
		r.Err = fmtError(errors.New("volume not exist when do remove"), q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	if err := volume.Remove(v, false); err != nil {
		r.Err = fmtError(err, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	delete(s.volumes, q.Name)
	assert(json.NewEncoder(resp).Encode(r))
}

func (s *stonePlugin) mount(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Mountpoint, Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received mount request for volume '%s'\n", q.Name)

	v, exists := s.volumes[q.Name]
	if !exists {
		r.Err = fmtError(errors.New("volume not exist when do mount"), q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	if _, err := v.Mount(); err != nil {
		r.Err = fmtError(err, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	m := v.GetPath()
	r.Mountpoint = &m
	assert(json.NewEncoder(resp).Encode(r))
}

func (s *stonePlugin) unmount(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received unmount request for volume '%s'\n", q.Name)

	v, exists := s.volumes[q.Name]
	if !exists {
		r.Err = fmtError(errors.New("volume not exist when do unmount"), q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	if err := v.Unmount(); err != nil {
		r.Err = fmtError(err, q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	assert(json.NewEncoder(resp).Encode(r))
}

func (s *stonePlugin) path(resp http.ResponseWriter, req *http.Request) {
	var q struct{ Name string }
	var r struct{ Mountpoint, Err *string }

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received path request for volume '%s'\n", q.Name)

	v, exists := s.volumes[q.Name]
	if !exists {
		r.Err = fmtError(errors.New("volume not exist when do path"), q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	m := v.GetPath()
	r.Mountpoint = &m
	assert(json.NewEncoder(resp).Encode(r))
}

func (s *stonePlugin) get(resp http.ResponseWriter, req *http.Request) {
	type Volume struct{ Name, Mountpoint string }

	var q struct{ Name string }
	var r struct {
		Volume *Volume
		Err    *string
	}

	assert(json.NewDecoder(req.Body).Decode(&q))
	log.Printf("Received get request for volume '%s'\n", q.Name)

	v, exists := s.volumes[q.Name]
	if !exists {
		r.Err = fmtError(errors.New("volume not exist when do get"), q.Name)
		assert(json.NewEncoder(resp).Encode(r))
		return
	}

	r.Volume = &Volume{
		Name:       q.Name,
		Mountpoint: v.GetPath(),
	}

	assert(json.NewEncoder(resp).Encode(r))
}

func (s *stonePlugin) list(resp http.ResponseWriter, req *http.Request) {
	type Volume struct{ Name, Mountpoint string }

	var r struct {
		Volumes []*Volume
		Err     *string
	}
	log.Println("Received list request")

	for _, vol := range s.volumes {
		r.Volumes = append(r.Volumes, &Volume{
			Name:       vol.GetName(),
			Mountpoint: vol.GetPath(),
		})
	}
	assert(json.NewEncoder(resp).Encode(r))
}

package mapper

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"rdcmd/internal/job"
	"strconv"
	"syscall"
	"time"
)

type IJobMapper interface {
	Add(job job.IJob) (int64, error)
	Delete(id int64) error
	Update(iJob job.IJob) error
	Get(id int64) (*job.IJob, error)
	List() ([]*job.IJob, error)

	SearchByName(name string) []*job.IJob
}

type FileJobMapper struct {
	path     string
	store    string
	delete   string
	lockFile *os.File
}

func NewFileJobMapper(path string) (*FileJobMapper, error) {
	store := path + "/store"
	del := path + "/delete"
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(store, os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(del, os.ModePerm); err != nil {
		return nil, err
	}
	return &FileJobMapper{
		path:   path,
		store:  store,
		delete: del,
	}, nil
}

func (l *FileJobMapper) lock() {
	for {
		f, err := os.Create(l.path + "/.lock")
		if err != nil {
			log.Debugf("failed to open lock, err:%+v", err)
			time.Sleep(20 * time.Microsecond)
			continue
		}
		l.lockFile = f
		err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err != nil {
			log.Debugf("failed to flock, err:%+v", err)
			time.Sleep(20 * time.Microsecond)
			continue
		}
		break
	}
}

func (l *FileJobMapper) unlock() {
	defer os.Remove(l.lockFile.Name())
	syscall.Flock(int(l.lockFile.Fd()), syscall.LOCK_UN)
}

func (l *FileJobMapper) getMaxId() (int64, error) {
	var maxId int64 = 0
	if fis, err := ioutil.ReadDir(l.store); err != nil {
		return -1, err
	} else {
		for i := range fis {
			fi := fis[i]
			if fi.IsDir() {
				name := fi.Name()
				if id, err := strconv.ParseInt(name, 10, 64); err != nil {
					log.Errorf("Invalid name:%s, err:%+v", name, err)
				} else {
					if maxId < id {
						maxId = id
					}
				}
			}
		}
	}
	return maxId + 1, nil
}

func (l *FileJobMapper) Add(job job.IJob) (int64, error) {
	l.lock()
	defer l.unlock()

	// generate job id
	id, err := l.getMaxId()
	if err != nil {
		return -1, err
	}
	job.SetId(id)

	// create job directory
	dir := fmt.Sprintf(l.store+"/%d", id)
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return -1, err
	}

	// write job entity json
	f, err := os.Create(dir + "/job.json")
	if err != nil {
		return -1, err
	}
	data, err := json.Marshal(job)
	if err != nil {
		return -1, err
	}
	if err = ioutil.WriteFile(f.Name(), data, os.ModePerm); err != nil {
		return -1, err
	}
	return id, nil
}

func (l *FileJobMapper) Delete(id int64) error {
	l.lock()
	defer l.unlock()

	// create job directory
	dir := fmt.Sprintf(l.store+"/%d", id)
	return os.RemoveAll(dir)
}

func (l *FileJobMapper) Update(iJob job.IJob) error {
	l.lock()
	defer l.unlock()

	if err := l.Delete(iJob.GetId()); err != nil {
		return err
	}

	// create job directory
	dir := fmt.Sprintf(l.store+"/%d", iJob.GetId())
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return err
	}

	// write job entity json
	f, err := os.Create(dir + "/job.json")
	if err != nil {
		return err
	}
	data, err := json.Marshal(iJob)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(f.Name(), data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (l *FileJobMapper) Get(id int64) (*job.IJob, error) {
	l.lock()
	defer l.unlock()

	fi := fmt.Sprintf(l.store+"/%d/job.json", id)
	bs, err := ioutil.ReadFile(fi)
	if err != nil && err != io.EOF {
		return nil, err
	}
	basicJob := job.BasicJob{}
	if err := json.Unmarshal(bs, &basicJob); err != nil {
		return nil, err
	}
	switch basicJob.GetType() {
	case job.SHELL:
		shellJob := job.ShellJob{}
		if err := json.Unmarshal(bs, &shellJob); err != nil {
			return nil, err
		}
		var j job.IJob = &shellJob
		return &j, nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid jobType:%+v", basicJob.GetType()))
	}
}

func (l *FileJobMapper) List() ([]*job.IJob, error) {
	l.lock()
	defer l.unlock()

	fis, err := ioutil.ReadDir(l.store)
	if err != nil {
		return nil, err
	}
	jobs := make([]*job.IJob, len(fis))
	for i := range fis {
		fi := fis[i]
		if fi.IsDir() {
			jobId, err := strconv.ParseInt(fi.Name(), 10, 64)
			if err != nil {
				log.Errorf("Invalid jobId: %+v, err:%+v", jobId, err)
			} else {
				j, err := l.Get(jobId)
				if err != nil {
					log.Errorf("Failed to get job, id:%d, err:%+v", jobId, err)
				} else {
					jobs[i] = j
				}
			}
		}
	}
	return jobs, nil
}

func (l *FileJobMapper) SearchByName(name string) []*job.IJob {
	panic("implement me")
}

package mapper

import (
	"rdcmd/internal/job"
	"testing"
)

func TestFileJobMapper_Add(t *testing.T) {
	fileMapper, err := NewFileJobMapper("/tmp/rdcmd")
	if err != nil {
		t.Error(err)
	}
	j, err := job.NewShellJob("myjob", "echo hello")
	if err != nil {
		t.Error(err)
	}
	id, err := fileMapper.Add(job.IJob(j))
	if err != nil {
		t.Error(err)
	}
	t.Logf("id:%+v", id)
}

func TestFileJobMapper_Delete(t *testing.T) {

}

func TestFileJobMapper_Get(t *testing.T) {
	fileMapper, err := NewFileJobMapper("/tmp/rdcmd")
	if err != nil {
		t.Error(err)
	}
	fis, err := fileMapper.List()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", fis)
}

func TestFileJobMapper_SearchByName(t *testing.T) {

}

func TestFileJobMapper_Update(t *testing.T) {

}

func TestFileJobMapper_getMaxId(t *testing.T) {

}

func TestFileJobMapper_lock(t *testing.T) {

}

func TestFileJobMapper_unlock(t *testing.T) {

}

func TestNewFileJobMapper(t *testing.T) {

}

package job

import (
	"fmt"
)

// IJob represent a job that can be executed and fetch result
type IJob interface {
	// GetId execute and return exit code
	GetId() int64
	SetId(id int64)

	GetInfo() string
	SetInfo(info string)
}

type Type int

const (
	SHELL Type = iota
)

type BasicJob struct {
	Type Type   `json:"type"`
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (b *BasicJob) GetType() Type {
	return b.Type
}

type ShellJob struct {
	BasicJob
	Cmd string `json:"cmd"`
}

func NewShellJob(name string, cmd string) (*ShellJob, error) {
	return &ShellJob{
		BasicJob: BasicJob{
			Name: name,
		},
		Cmd: cmd,
	}, nil
}

func (s *ShellJob) GetId() int64 {
	return s.Id
}

func (s *ShellJob) GetInfo() string {
	return fmt.Sprintf("SHELLJOB, id:%d, name:%s\n>%s",
		s.Id, s.Name, s.Cmd)
}

func (s *ShellJob) SetId(id int64) {
	s.Id = id
}

func (s *ShellJob) SetInfo(info string) {
}

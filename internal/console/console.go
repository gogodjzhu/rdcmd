package console

import "rdcmd/internal/job"

type IConsole interface {
	UpdateJob(job *job.IJob) error
	DeleteJob(iJob *job.IJob) error
	ListJobs() []*job.IJob
	ExecuteJob(iJob *job.IJob) error
}

type DefaultConsole struct {
}

func (d DefaultConsole) UpdateJob(job *job.IJob) error {
	panic("implement me")
}

func (d DefaultConsole) DeleteJob(iJob *job.IJob) error {
	panic("implement me")
}

func (d DefaultConsole) ListJobs() []*job.IJob {
	panic("implement me")
}

func (d DefaultConsole) ExecuteJob(iJob *job.IJob) error {
	panic("implement me")
}

package worker

import (
	"context"
	"website-monitor/internal/checker"
)

type Job struct {
	SiteID int
	URL    string
}

type Worker struct {
	jobs    <-chan Job                 //worker get only tasks
	results chan<- checker.CheckResult //worker only send results
}

func NewWorker(
	jobs <-chan Job,
	results chan<- checker.CheckResult,
) *Worker {
	return &Worker{
		jobs:    jobs,
		results: results,
	}
}

func (w *Worker) Start(ctx context.Context) {

}

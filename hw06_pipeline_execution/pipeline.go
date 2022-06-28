package hw06pipelineexecution

import (
	"sort"
	"sync"
)

type (
	In       = <-chan interface{}
	Out      = In
	Bi       = chan interface{}
	St       = chan StageRes
	StageRes struct {
		WorkerIndex int
		Value       interface{}
	}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	data := []interface{}{}
	for val := range in {
		data = append(data, val)
	}
	isDone := false
	mutex := sync.Mutex{}

	go func() {
		<-done
		mutex.Lock()
		isDone = true
		mutex.Unlock()
	}()

	wg := sync.WaitGroup{}
	wg.Add(len(data))
	stageResCh := make(St, len(data))

	for i, val := range data {
		go func(workerIndex int, val interface{}) {
			defer wg.Done()
			for _, stage := range stages {
				stage := stage
				mutex.Lock()
				if isDone {
					mutex.Unlock()
					return
				}
				mutex.Unlock()

				out := make(Bi, 1)
				defer close(out)

				out <- val
				val = <-stage(out)
			}

			stageResCh <- StageRes{WorkerIndex: workerIndex, Value: val}
		}(i, val)
	}
	wg.Wait()
	close(stageResCh)

	stagingResults := []StageRes{}
	for r := range stageResCh {
		stagingResults = append(stagingResults, r)
	}

	sort.Slice(stagingResults, func(i, j int) bool {
		return stagingResults[i].WorkerIndex < stagingResults[j].WorkerIndex
	})

	resCh := make(Bi, len(data))
	defer close(resCh)

	for _, r := range stagingResults {
		resCh <- r.Value
	}

	return resCh
}

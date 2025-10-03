package types

import (
	"sync"
	"time"
)

/*
The idea behind this is that every time a Model Provider is intialized
We can retrieve the model list and stash here. Then the ModelProvider can check to see if
the local cache is stale, in which case it can pull it fresh from the API.
*/

type Model struct {
	Name       string   `json:"name"`
	Modalities []string `json:"modalities"`
}

type ModelRefresher struct {
	Expiry      time.Duration
	Models      []Model
	Mutex       sync.RWMutex
	LastUpdated time.Time
}

func NewModelRefresher(refreshTimeMS int64) *ModelRefresher {
	return &ModelRefresher{
		Expiry: time.Duration(refreshTimeMS) * time.Second,
		Models: []Model{},
		Mutex:  sync.RWMutex{},
	}
}

func (mf *ModelRefresher) StashModels(models []Model) error {
	mf.Mutex.Lock()
	mf.Models = models
	mf.Mutex.Unlock()
	mf.LastUpdated = time.Now()
	return nil
}

func (mf *ModelRefresher) RetrieveModels() []Model {
	modelsCopy := []Model{}
	mf.Mutex.RLock()
	modelsCopy = append(modelsCopy, mf.Models...)
	mf.Mutex.RUnlock()
	return modelsCopy
}

func (mf *ModelRefresher) IsStale() bool {
	return time.Since(mf.LastUpdated) > mf.Expiry
}

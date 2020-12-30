package destinations

import (
	"github.com/jitsucom/eventnative/events"
	"github.com/jitsucom/eventnative/storages"
)

//map["tokenId"]map["destination_name"]true
type TokenizedIds map[string]map[string]bool

//map["tokenId"]map["destination_name"]interface
//because 1 token id = ∞ storages
type TokenizedStorages map[string]map[string]storages.StorageProxy

//map["tokenId"]map["tokenId | destination_name"]interface
//because 1 token id = 1logger but ∞ event.queue
type TokenizedConsumers map[string]map[string]events.Consumer

func (ts TokenizedStorages) Add(tokenId, name string, proxy storages.StorageProxy) {
	storageProxies, ok := ts[tokenId]
	if !ok {
		storageProxies = map[string]storages.StorageProxy{}
		ts[tokenId] = storageProxies
	}
	storageProxies[name] = proxy
}

func (ts TokenizedStorages) AddAll(other TokenizedStorages) {
	for tokenId, storages := range other {
		for name, storage := range storages {
			ts.Add(tokenId, name, storage)
		}
	}
}

func (tc TokenizedConsumers) Add(tokenId, name string, proxy events.Consumer) {
	consumers, ok := tc[tokenId]
	if !ok {
		consumers = map[string]events.Consumer{}
		tc[tokenId] = consumers
	}
	consumers[name] = proxy
}

func (tc TokenizedConsumers) AddAll(other TokenizedConsumers) {
	for tokenId, consumers := range other {
		for name, consumer := range consumers {
			tc.Add(tokenId, name, consumer)
		}
	}
}

func (ti TokenizedIds) Add(tokenId, name string) {
	ids, ok := ti[tokenId]
	if !ok {
		ids = map[string]bool{}
		ti[tokenId] = ids
	}
	ids[name] = true
}

func (ti TokenizedIds) AddAll(other TokenizedIds) {
	for tokenId, ids := range other {
		for id := range ids {
			ti.Add(tokenId, id)
		}
	}
}

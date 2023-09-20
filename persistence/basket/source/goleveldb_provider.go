package source

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"
)

type GoLevelDBProvider struct {
	snapshotInterval int
	mu               sync.RWMutex
	db               *leveldb.DB
}

func NewGoLevelDBProvider(snapshotInterval int, dir string) (*GoLevelDBProvider, error) {
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	return &GoLevelDBProvider{
		snapshotInterval: snapshotInterval,
		db:               db,
	}, nil
}

// DeleteEvents removes all events from the provider
func (provider *GoLevelDBProvider) DeleteEvents(actorName string, inclusiveToIndex int) {
}

func (provider *GoLevelDBProvider) Restart() {}

func (provider *GoLevelDBProvider) GetSnapshotInterval() int {
	return provider.snapshotInterval
}

func (provider *GoLevelDBProvider) GetSnapshot(actorName string) (snapshot interface{}, eventIndex int, ok bool) {
	maxEventIndex := -1
	prefix := fmt.Sprintf("snapshot:%s:", actorName)
	iter := provider.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()
	for iter.Next() {
		en := &envelope{}
		if err := json.Unmarshal(iter.Value(), en); err != nil {
			log.Fatal(err)
		}
		return en.message(), maxEventIndex, true
	}
	if err := iter.Error(); err != nil {
		return nil, 0, false
	}
	return nil, maxEventIndex, false
}

func (provider *GoLevelDBProvider) PersistSnapshot(actorName string, eventIndex int, snapshot proto.Message) {
	key := fmt.Sprintf("snapshot:%s:%s", actorName, strconv.Itoa(eventIndex))
	envelope := newEnvelope(snapshot, "snapshot", eventIndex)
	m, err := json.Marshal(envelope)
	if err != nil {
		log.Fatal(err)
	}
	_ = provider.db.Put([]byte(key), m, nil)
}

func (provider *GoLevelDBProvider) DeleteSnapshots(actorName string, inclusiveToIndex int) {
}

func (provider *GoLevelDBProvider) PersistEvent(actorName string, eventIndex int, snapshot proto.Message) {
	key := fmt.Sprintf("event:%s:%s", actorName, strconv.Itoa(eventIndex))
	envelope := newEnvelope(snapshot, "event", eventIndex)
	m, err := json.Marshal(envelope)
	if err != nil {
		log.Fatal(err)
	}
	_ = provider.db.Put([]byte(key), m, nil)
}

func (provider *GoLevelDBProvider) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{})) {
	for i := eventIndexStart; i <= eventIndexEnd; i++ {
		key := fmt.Sprintf("event:%s:%s", actorName, strconv.Itoa(i))
		r, _ := provider.db.Get([]byte(key), nil)
		en := &envelope{}
		if err := json.Unmarshal(r, en); err != nil {
			log.Fatal(err)
		}
		callback(en.message())
	}
}

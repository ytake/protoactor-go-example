package source

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/asynkron/protoactor-go/persistence"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"
)

type GoLevelDBProvider struct {
	snapshotInterval int
	db               *leveldb.DB
	dbDir            string
	mu               sync.RWMutex
}

func NewGoLevelDBProvider(snapshotInterval int, dbDir string) (*GoLevelDBProvider, error) {
	db, err := leveldb.OpenFile(dbDir, nil)
	if err != nil {
		return nil, err
	}
	return &GoLevelDBProvider{
		snapshotInterval: snapshotInterval,
		db:               db,
		dbDir:            dbDir,
		mu:               sync.RWMutex{},
	}, nil
}

// DeleteEvents removes all events from the provider
func (provider *GoLevelDBProvider) DeleteEvents(_ string, _ int) {
}

func (provider *GoLevelDBProvider) Restart() {
	_ = provider.db.Close()
	db, err := leveldb.OpenFile(provider.dbDir, nil)
	if err != nil {
		log.Fatal(err)
	}
	provider.db = db
}

func (provider *GoLevelDBProvider) GetSnapshotInterval() int {
	return provider.snapshotInterval
}

func (provider *GoLevelDBProvider) GetSnapshot(actorName string) (snapshot interface{}, eventIndex int, ok bool) {
	prefix := fmt.Sprintf("snapshot:%s:", actorName)
	iter := provider.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()
	for iter.Next() {
		if iter.Last() {
			en := &envelope{}
			if err := json.Unmarshal(iter.Value(), en); err != nil {
				log.Fatal(err)
			}
			return en.message(), en.EventIndex, true
		}
	}

	if err := iter.Error(); err != nil {
		return nil, 0, false
	}
	return nil, 0, false
}

func (provider *GoLevelDBProvider) PersistSnapshot(actorName string, eventIndex int, snapshot proto.Message) {
	provider.mu.Lock()
	defer provider.mu.Unlock()
	key := fmt.Sprintf("snapshot:%s:%s", actorName, strconv.Itoa(eventIndex))
	envelope := newEnvelope(snapshot, "snapshot", eventIndex)
	m, err := json.Marshal(envelope)
	if err != nil {
		log.Fatal(err)
	}
	_ = provider.db.Put([]byte(key), m, nil)
}

func (provider *GoLevelDBProvider) DeleteSnapshots(_ string, _ int) {
}

func (provider *GoLevelDBProvider) PersistEvent(actorName string, eventIndex int, snapshot proto.Message) {
	provider.mu.Lock()
	defer provider.mu.Unlock()
	key := fmt.Sprintf("event:%s:%s", actorName, strconv.Itoa(eventIndex))
	envelope := newEnvelope(snapshot, "event", eventIndex)
	m, err := json.Marshal(envelope)
	if err != nil {
		log.Fatal(err)
	}
	_ = provider.db.Put([]byte(key), m, nil)
}

func (provider *GoLevelDBProvider) GetEvents(actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{})) {
	prefix := fmt.Sprintf("event:%s:", actorName)
	iter := provider.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()
	for iter.Next() {
		parts := strings.Split(string(iter.Key()), ":")
		i, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Fatal(err)
		}
		// 指定のイベントの開始位置から再読み込みをおこなう
		if i >= eventIndexStart {
			en := &envelope{}
			if err = json.Unmarshal(iter.Value(), en); err != nil {
				log.Fatal(err)
			}
			callback(en.message())
		}
	}
}

func (provider *GoLevelDBProvider) GetState() persistence.ProviderState {
	return provider
}

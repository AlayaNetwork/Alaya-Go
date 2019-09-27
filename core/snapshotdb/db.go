package snapshotdb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"path"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbError "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/journal"
	"github.com/syndtr/goleveldb/leveldb/memdb"
)

func getBaseDBPath(dbpath string) string {
	return path.Join(dbpath, DBBasePath)
}

func newDB(stor storage) (*snapshotDB, error) {
	dbpath := stor.Path()
	baseDB, err := leveldb.OpenFile(getBaseDBPath(dbpath), nil)
	if err != nil {
		return nil, fmt.Errorf("[SnapshotDB]open baseDB fail:%v", err)
	}
	unCommitBlock := new(unCommitBlocks)
	unCommitBlock.blocks = make(map[common.Hash]*blockData)
	return &snapshotDB{
		path:            dbpath,
		storage:         stor,
		unCommit:        unCommitBlock,
		committed:       make([]*blockData, 0),
		baseDB:          baseDB,
		current:         newCurrent(dbpath),
		snapshotLockC:   snapshotUnLock,
		exitCh:          make(chan struct{}, 0),
		currentUpdateCh: make(chan struct{}, 1),
	}, nil
}

func (s *snapshotDB) getBlockFromJournal(fd fileDesc) (*blockData, error) {
	reader, err := s.storage.Open(fd)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	journals := journal.NewReader(reader, nil, false, false)
	j, err := journals.Next()
	if err != nil {
		return nil, err
	}
	var header journalHeader
	if err := decode(j, &header); err != nil {
		return nil, err
	}
	block := new(blockData)
	block.ParentHash = header.ParentHash
	block.kvHash = header.KvHash
	if fd.BlockHash != s.getUnRecognizedHash() {
		block.BlockHash = fd.BlockHash
	}
	block.Number = new(big.Int).SetUint64(fd.Num)
	block.data = memdb.New(DefaultComparer, 0)
	block.readOnly = true
	/*switch header.From {
	case journalHeaderFromUnRecognized:
		if fd.BlockHash == s.getUnRecognizedHash() {
			block.readOnly = false
		} else {
			block.readOnly = true
		}
	case journalHeaderFromRecognized:
		if fd.Num <= s.current.HighestNum.Uint64() {
			block.readOnly = true
		}
	}*/
	//var kvhash common.Hash
	for {
		j, err := journals.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		var body journalData
		if err := decode(j, &body); err != nil {
			return nil, err
		}
		if err := block.data.Put(body.Key, body.Value); err != nil {
			return nil, err
		}
		//kvhash = body.Hash
	}
	//block.kvHash = kvhash
	return block, nil
}

func (s *snapshotDB) recover(stor storage) error {
	dbpath := stor.Path()
	c, err := loadCurrent(dbpath)
	if err != nil {
		return fmt.Errorf("[SnapshotDB.recover]load  current fail:%v", err)
	}
	s.path = dbpath

	if blockchain != nil {
		currentHead := blockchain.CurrentHeader()
		if c.HighestNum.Cmp(currentHead.Number) != 0 {
			s.current = c
			if err := s.SetCurrent(currentHead.Hash(), *c.BaseNum, *currentHead.Number); err != nil {
				return err
			}
		} else {
			s.current = c
		}
	} else {
		s.current = c
	}

	//baseDB
	baseDB, err := leveldb.OpenFile(getBaseDBPath(dbpath), nil)
	if _, corrupted := err.(*leveldbError.ErrCorrupted); corrupted {
		baseDB, err = leveldb.RecoverFile(getBaseDBPath(dbpath), nil)
		if err != nil {
			return fmt.Errorf("[SnapshotDB.recover]RecoverFile baseDB fail:%v", err)
		}
	}
	if err != nil {
		return fmt.Errorf("[SnapshotDB.recover]open baseDB fail:%v", err)
	}
	s.baseDB = baseDB

	//storage
	s.storage = stor
	fds, err := s.storage.List(TypeJournal)
	sortFds(fds)
	baseNum := s.current.BaseNum.Uint64()
	highestNum := s.current.HighestNum.Uint64()
	//	UnRecognizedHash := s.getUnRecognizedHash()
	s.committed = make([]*blockData, 0)
	//	s.journalw = make(map[common.Hash]*journalWriter)
	unCommitBlock := new(unCommitBlocks)
	unCommitBlock.blocks = make(map[common.Hash]*blockData)
	s.unCommit = unCommitBlock
	s.snapshotLockC = snapshotUnLock
	s.exitCh = make(chan struct{}, 0)
	s.currentUpdateCh = make(chan struct{}, 1)

	//read Journal
	for _, fd := range fds {
		if baseNum < fd.Num && fd.Num <= highestNum {
			if header := blockchain.GetHeaderByHash(fd.BlockHash); header != nil {
				block, err := s.getBlockFromJournal(fd)
				if err != nil {
					return err
				}
				s.committed = append(s.committed, block)
				logger.Debug("recover block ", "num", block.Number, "hash", block.BlockHash.String())
				continue
			}
		}
		if err := s.storage.Remove(fd); err != nil {
			return err
		}
	}
	return nil
}

/*
func (s *snapshotDB) rmOldRecognizedBlockData() error {
	s.unCommit.Lock()
	defer s.unCommit.Unlock()
	for key, value := range s.unCommit.blocks {
		if s.current.HighestNum.Cmp(value.Number) >= 0 {
			if !value.readOnly {
				if err := value.jWriter.Close(); err != nil {
					return err
				}
			}
			delete(s.unCommit.blocks, key)
			if err := s.rmJournalFile(value.Number, key); err != nil {
				return err
			}
		}
	}
	return nil
}
*/
func (s *snapshotDB) generateKVHash(k, v []byte, hash common.Hash) common.Hash {
	var buf bytes.Buffer
	buf.Write(k)
	buf.Write(v)
	buf.Write(hash.Bytes())
	return rlpHash(buf.Bytes())
}

func (s *snapshotDB) getUnRecognizedHash() common.Hash {
	return common.ZeroHash
}

//
//func (s *snapshotDB) closeJournalWriter(hash common.Hash) error {
//	s.journalWriterLock.Lock()
//	defer s.journalWriterLock.Unlock()
//	if j, ok := s.journalw[hash]; ok {
//		if err := j.Close(); err != nil {
//			return errors.New("[snapshotdb]close  journal writer fail:" + err.Error())
//		}
//		delete(s.journalw, hash)
//	}
//	return nil
//}
//
//const (
//	hashLocationUnCommitted = 2
//	hashLocationCommitted   = 3
//	hashLocationNotFound    = 4
//)

//
//func (s *snapshotDB) checkHashChain(hash common.Hash) (int, bool) {
//	var (
//		lastBlockNumber = big.NewInt(0)
//		lastParentHash  = hash
//	)
//	// find from unCommit
//	for {
//		if block, ok := s.unCommit.Load(lastParentHash); ok {
//			if lastParentHash == block.ParentHash {
//				logger.Error("loop error")
//				return 0, false
//			}
//			lastParentHash = block.ParentHash
//			lastBlockNumber = block.Number
//		} else {
//			break
//		}
//	}
//
//	//check  last block find from unCommit is right
//	if lastBlockNumber.Int64() > 0 {
//		if s.current.HighestNum.Int64() != lastBlockNumber.Int64()-1 {
//			logger.Error("[snapshotDB] find lastblock  fail ,num not compare", "current", s.current.HighestNum, "last", lastBlockNumber.Int64()-1)
//			return 0, false
//		}
//		if s.current.HighestHash == common.ZeroHash {
//			return hashLocationUnCommitted, true
//		}
//		if s.current.HighestHash != lastParentHash {
//			logger.Error("[snapshotDB] find lastblock  fail ,hash not compare", "current", s.current.HighestHash.String(), "last", lastParentHash.String())
//			return 0, false
//		}
//		return hashLocationUnCommitted, true
//	}
//	// if not find from unCommit, find from committed
//	for _, value := range s.committed {
//		if value.BlockHash == hash {
//			return hashLocationCommitted, true
//		}
//	}
//	return hashLocationNotFound, true
//}

func (s *snapshotDB) put(hash common.Hash, key, value []byte) error {
	s.unCommit.Lock()
	defer s.unCommit.Unlock()
	block, ok := s.unCommit.blocks[hash]
	if !ok {
		return fmt.Errorf("not find the block by hash:%v", hash.String())
	}
	if block.readOnly {
		return errors.New("can't put read only block")
	}
	block.kvHash = s.generateKVHash(key, value, block.kvHash)
	if err := block.data.Put(key, value); err != nil {
		return err
	}
	return nil
}

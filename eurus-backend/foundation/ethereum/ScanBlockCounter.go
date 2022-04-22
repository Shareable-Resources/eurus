package ethereum

import (
	"eurus-backend/foundation"
	"eurus-backend/foundation/database"
	"eurus-backend/foundation/log"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ScanBlockCounterMode int

const (
	ScanBlockModeContinuous ScanBlockCounterMode = iota
	ScanBlockModeDiscret
)

type ScanBlockStatus int

const (
	ScanBlockStatusError ScanBlockStatus = iota - 1
	ScanBlockStatusOnGoing
	ScanBlockStatusIncomplete
	ScanBlockActionFinished
)

type RescanBlock struct {
	database.DbModel
	EventId    decimal.Decimal `gorm:"primaryKey;autoIncrement"`
	ChainId    uint64
	StartNum   decimal.Decimal
	EndNum     *decimal.Decimal
	CurrentNum *decimal.Decimal
	Status     ScanBlockStatus
}

type ScanBlockCounter struct {
	db            *database.Database
	mode          ScanBlockCounterMode
	RescanHandler func(sender *ScanBlockCounter, eventId decimal.Decimal, chainId uint64, from decimal.Decimal, to decimal.Decimal)
	chainId       uint64
	one           *big.Int
	decimalOne    decimal.Decimal
	cri           *sync.Mutex
	LoggerName    string
}

func NewScanBlockCounter(mode ScanBlockCounterMode, chainId uint64,
	rescanHandler func(sender *ScanBlockCounter, eventId decimal.Decimal, chainId uint64, from decimal.Decimal, to decimal.Decimal), dbFileName string, LoggerName string) (*ScanBlockCounter, error) {
	if len(dbFileName) == 0 {
		dbFileName = "rescan.db"
	}
	counter := new(ScanBlockCounter)
	counter.mode = mode
	counter.db = new(database.Database)
	counter.db.DBName = dbFileName
	counter.db.DbType = database.DatabaseSqlite
	counter.db.Logger = log.GetLogger(log.Name.Root)
	counter.RescanHandler = rescanHandler
	counter.chainId = chainId
	counter.one = big.NewInt(1)
	counter.decimalOne = decimal.NewFromBigInt(counter.one, 0)
	counter.cri = new(sync.Mutex)
	counter.LoggerName = LoggerName
	err := counter.initSqlite()
	return counter, err
}

func (me *ScanBlockCounter) initSqlite() error {
	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}
	tx := conn.Session(&gorm.Session{}).Exec("CREATE TABLE IF NOT EXISTS rescan_blocks (event_id INTEGER PRIMARY KEY AUTOINCREMENT, chain_id INTEGER NOT NULL, start_num INTEGER NOT NULL, end_num INTEGER, current_num INTEGER, status INTEGER NOT NULL, created_date DATETIME NOT NULL, last_modified_date DATETIME NOT NULL);")
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Create table error")
	}

	tx = conn.Session(&gorm.Session{}).Exec("CREATE index if not exists rescan_blocks_idx1 ON rescan_blocks (status, chain_id);")
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Create idx1 error")
	}

	tx = conn.Session(&gorm.Session{}).Exec("CREATE index if not exists rescan_blocks_idx2 ON rescan_blocks (event_id, chain_id) ;")
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Create idx2 error")
	}

	conn.Session(&gorm.Session{}).Exec("VACUUM")
	return nil
}

func (me *ScanBlockCounter) Start() error {
	var rescanBlockList []*RescanBlock = make([]*RescanBlock, 0)

	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	dbTx := conn.Where("status = ? AND chain_id = ?", ScanBlockStatusIncomplete, me.chainId).Order("start_num").Find(&rescanBlockList)
	if dbTx.Error != nil {
		return dbTx.Error
	}

	if len(rescanBlockList) > 0 {
		go func() {
			for _, rescanBlock := range rescanBlockList {
				var startFrom decimal.Decimal

				if rescanBlock.CurrentNum != nil {
					oneDec := decimal.NewFromInt(1)
					startFrom = rescanBlock.CurrentNum.Add(oneDec)
				} else {
					startFrom = rescanBlock.StartNum
				}
				me.RescanHandler(me, rescanBlock.EventId, rescanBlock.ChainId, startFrom, *rescanBlock.EndNum)
			}
		}()
	}
	return nil
}

func (me *ScanBlockCounter) UpdateLatestBlock(blockNum *big.Int) error {
	me.cri.Lock()
	defer me.cri.Unlock()

	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	rescanBlock := new(RescanBlock)
	dbTx := conn.Where("status = ? and chain_id = ?", ScanBlockStatusOnGoing, me.chainId).Order("start_num DESC").Find(&rescanBlock)
	if dbTx.Error != nil {
		return dbTx.Error
	}

	if dbTx.RowsAffected == 0 {
		//Newly created case
		rescanBlock.InitDate()
		rescanBlock.StartNum = decimal.NewFromBigInt(blockNum, 0)
		rescanBlock.Status = ScanBlockStatusOnGoing
		rescanBlock.ChainId = me.chainId
		conn.Create(rescanBlock)
	} else {
		if me.mode == ScanBlockModeContinuous {
			delta := big.NewInt(0)
			delta = delta.Sub(blockNum, rescanBlock.StartNum.BigInt())
			cmpResult := delta.Cmp(me.one)

			if cmpResult == 0 {
				//Scan block no gap
				rescanBlock.StartNum = decimal.NewFromBigInt(blockNum, 0)
				rescanBlock.LastModifiedDate = time.Now()
				rescanBlock.Status = ScanBlockStatusOnGoing

				dbTx := conn.Model(rescanBlock).Select("start_num", "last_modified_time").Updates(*rescanBlock)
				if dbTx.Error != nil {
					return dbTx.Error
				}
				//log.GetLogger(log.Name.Root).Debugln("Start processing block number: ", blockNum.String())
			} else if cmpResult > 0 {
				//Scan block gap found
				endBlock := big.NewInt(0)
				endBlock = endBlock.Sub(blockNum, me.one)

				endNum := decimal.NewFromBigInt(endBlock, 0)
				rescanBlock.EndNum = &endNum
				rescanBlock.LastModifiedDate = time.Now()
				rescanBlock.Status = ScanBlockStatusIncomplete

				newOnGoingBlock := new(RescanBlock)
				newOnGoingBlock.StartNum = decimal.NewFromBigInt(blockNum, 0)
				newOnGoingBlock.Status = ScanBlockStatusOnGoing
				newOnGoingBlock.ChainId = me.chainId
				newOnGoingBlock.InitDate()

				err := conn.Transaction(func(dbTx *gorm.DB) error {
					dbTx1 := dbTx.Session(&gorm.Session{}).Model(rescanBlock).Select("start_num, end_num, status, last_modified_date").Updates(*rescanBlock)
					if dbTx1.Error != nil {
						return dbTx1.Error
					}
					dbTx1 = dbTx.Session(&gorm.Session{}).Create(newOnGoingBlock)
					return dbTx1.Error
				})

				if err != nil {
					return err
				}
				startNum := rescanBlock.StartNum.Add(me.decimalOne)
				log.GetLogger(me.LoggerName).Infoln("Lastest block gap found. Trigger rescan. Latest block received: ", blockNum.String(), " event id: ", rescanBlock.EventId.String(), " start num: ", startNum.String(), " end num: ", rescanBlock.EndNum.String())
				go me.RescanHandler(me, rescanBlock.EventId, rescanBlock.ChainId, startNum, *rescanBlock.EndNum)

			} else {
				//error case, input block number < last block number
				return foundation.NewErrorWithMessage(foundation.InvalidArgument, "Block number is too small. Latest block number:  "+rescanBlock.StartNum.BigInt().String())
			}
		} else if me.mode == ScanBlockModeDiscret {
			cmpResult := blockNum.Cmp(rescanBlock.StartNum.BigInt())
			if cmpResult < 0 {
				return foundation.NewErrorWithMessage(foundation.InvalidArgument, "Block number is too small. Latest block number:  "+rescanBlock.StartNum.BigInt().String())
			} else {
				//If ScanBlockModeDiscret, always update record
				rescanBlock.StartNum = decimal.NewFromBigInt(blockNum, 0)
				rescanBlock.LastModifiedDate = time.Now()
				rescanBlock.Status = ScanBlockStatusOnGoing

				dbTx := conn.Model(rescanBlock).Select("start_num", "last_modified_time").Updates(*rescanBlock)
				if dbTx.Error != nil {
					return dbTx.Error
				}
				//log.GetLogger(log.Name.Root).Debugln("Start processing block number: ", blockNum.String())
			}
		} else {
			return foundation.NewErrorWithMessage(foundation.InvalidArgument, "Invalid mode")
		}
	}
	return nil
}

func (me *ScanBlockCounter) UpdateRescanEvent(eventId decimal.Decimal, blockNum *big.Int) (ScanBlockStatus, error) {
	me.cri.Lock()
	defer me.cri.Unlock()
	conn, err := me.db.GetConn()

	if err != nil {
		return ScanBlockStatusError, err
	}

	// defer func() {
	// 	sqlDb, err1 := conn.DB()
	// 	if err1 == nil {
	// 		sqlDb.Close()
	// 	}
	// }()

	rescanBlock := new(RescanBlock)
	dbTx := conn.Session(&gorm.Session{}).Where("event_id = ?", eventId).First(rescanBlock)
	if dbTx.Error != nil {
		return ScanBlockStatusError, err
	}

	if dbTx.RowsAffected == 0 {
		return ScanBlockStatusError, foundation.NewErrorWithMessage(foundation.RecordNotFound, "Event not found")
	}

	if rescanBlock.Status != ScanBlockStatusIncomplete {
		return ScanBlockStatusError, foundation.NewErrorWithMessage(foundation.InvalidArgument, "Event status is not incomplete")
	}

	if rescanBlock.CurrentNum == nil {
		//Newly rescan case
		if me.mode == ScanBlockModeContinuous {
			resultDelta := big.NewInt(0)
			resultDelta = resultDelta.Sub(blockNum, rescanBlock.StartNum.BigInt())
			if resultDelta.Cmp(me.one) != 0 {
				return ScanBlockStatusError, foundation.NewErrorWithMessage(foundation.InvalidArgument, "Input block number is not next to the start number: "+rescanBlock.StartNum.BigInt().String())
			}
		}
		return me.dbIncreaseCurrentBlockCounter(conn, blockNum, rescanBlock)
	}

	//Re-scanning case
	delta := big.NewInt(0)

	if me.mode == ScanBlockModeContinuous {
		delta = delta.Sub(blockNum, rescanBlock.CurrentNum.BigInt())
		if delta.Cmp(me.one) == 0 {
			return me.dbIncreaseCurrentBlockCounter(conn, blockNum, rescanBlock)
		}
	} else {
		if blockNum.Cmp(rescanBlock.CurrentNum.BigInt()) >= 0 {
			return me.dbIncreaseCurrentBlockCounter(conn, blockNum, rescanBlock)
		} else {
			return ScanBlockStatusError, foundation.NewErrorWithMessage(foundation.InvalidArgument, "Input block number smaller than or equal to the last block number. Last block number: "+rescanBlock.CurrentNum.BigInt().String())
		}

	}
	return ScanBlockStatusError, foundation.NewErrorWithMessage(foundation.InvalidArgument, "Input block number is not next to the last scan block number")
}

func (me *ScanBlockCounter) TriggerNewRescanEvent(lastestBlockNum *big.Int) error {
	endBlock := lastestBlockNum
	rescanBlock, err := me.GetLatestProgress()
	if err != nil {
		return err
	}

	me.cri.Lock()
	defer me.cri.Unlock()

	endNum := decimal.NewFromBigInt(endBlock, 0)
	endNum = endNum.Sub(decimal.NewFromInt(1))
	rescanBlock.EndNum = &endNum
	rescanBlock.LastModifiedDate = time.Now()
	rescanBlock.Status = ScanBlockStatusIncomplete

	newOnGoingBlock := new(RescanBlock)

	nextBlockNum := big.NewInt(0)
	nextBlockNum = nextBlockNum.Add(lastestBlockNum, me.one)

	newOnGoingBlock.StartNum = decimal.NewFromBigInt(nextBlockNum, 0)
	newOnGoingBlock.Status = ScanBlockStatusOnGoing
	newOnGoingBlock.ChainId = me.chainId
	newOnGoingBlock.InitDate()

	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	err = conn.Transaction(func(dbTx *gorm.DB) error {
		dbTx = dbTx.Model(rescanBlock).Select("start_num, end_num, status, last_modified_date").Updates(*rescanBlock)
		if dbTx.Error != nil {
			return dbTx.Error
		}
		dbTx = dbTx.Create(newOnGoingBlock)
		return dbTx.Error
	})

	if err != nil {
		return err
	}

	go me.RescanHandler(me, rescanBlock.EventId, rescanBlock.ChainId, rescanBlock.StartNum, *rescanBlock.EndNum)
	return nil
}

func (me *ScanBlockCounter) CreateNewRescanInRange(from *big.Int, to *big.Int) error {

	me.cri.Lock()
	defer me.cri.Unlock()

	conn, err := me.db.GetConn()
	if err != nil {
		return err
	}

	newRescan := new(RescanBlock)
	newRescan.InitDate()
	newRescan.StartNum = decimal.NewFromBigInt(from, 0)
	endNum := decimal.NewFromBigInt(to, 0)
	newRescan.EndNum = &endNum
	newRescan.ChainId = me.chainId
	newRescan.Status = ScanBlockStatusIncomplete

	dbTx := conn.Create(newRescan)
	if dbTx.Error != nil {
		return dbTx.Error
	}

	go me.RescanHandler(me, newRescan.EventId, newRescan.ChainId, newRescan.StartNum, *newRescan.EndNum)
	return nil
}

func (me *ScanBlockCounter) GetLatestProgress() (*RescanBlock, error) {
	me.cri.Lock()
	defer me.cri.Unlock()

	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}
	lastestProgress := new(RescanBlock)

	dbTx := conn.Where("status = ? AND chain_id = ?", ScanBlockStatusOnGoing, me.chainId).First(&lastestProgress)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	return lastestProgress, nil
}

func (me *ScanBlockCounter) GetEventProgress(eventId decimal.Decimal) (*RescanBlock, error) {
	me.cri.Lock()
	defer me.cri.Unlock()

	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}
	progress := new(RescanBlock)

	dbTx := conn.Where("event_id = ? AND chain_id = ?", eventId, me.chainId).First(&progress)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}

	return progress, nil
}

func (me *ScanBlockCounter) GetAllEventProgress() ([]*RescanBlock, error) {
	me.cri.Lock()
	defer me.cri.Unlock()
	conn, err := me.db.GetConn()
	if err != nil {
		return nil, err
	}

	var rescanBlockList []*RescanBlock = make([]*RescanBlock, 0)
	dbTx := conn.Where("chain_id = ?", me.chainId).Find(&rescanBlockList)

	return rescanBlockList, dbTx.Error
}

func (me *ScanBlockCounter) dbIncreaseCurrentBlockCounter(conn *gorm.DB, blockNum *big.Int, rescanBlock *RescanBlock) (ScanBlockStatus, error) {

	currNum := decimal.NewFromBigInt(blockNum, 0)
	rescanBlock.CurrentNum = &currNum
	rescanBlock.LastModifiedDate = time.Now()

	if rescanBlock.EndNum.BigInt().Cmp(blockNum) <= 0 {
		//Rescan procedure completed
		rescanBlock.Status = ScanBlockActionFinished
	}

	dbTx := conn.Session(&gorm.Session{}).Model(rescanBlock).Select("current_num, status, last_modified_date").Updates(*rescanBlock)
	if dbTx.Error != nil {
		return ScanBlockStatusError, dbTx.Error
	}

	return rescanBlock.Status, nil
}

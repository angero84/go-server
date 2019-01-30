package kdatabase

import (

	"database/sql"
	"github.com/angero84/go-server/kobject"
	"errors"

	klog "github.com/angero84/go-server/klogger"
	"github.com/angero84/go-server/kutil"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type KDB struct {
	*kobject.KObject
	db			*sql.DB
	info		*KDBInfo
	connOpt		*KDBConnOpt
}

func NewKDB(info *KDBInfo, connOpt *KDBConnOpt) (object *KDB, err error){

	if nil == info {
		err = errors.New("NewKDB() info is nil")
		return
	}

	if nil == connOpt {
		connOpt = &KDBConnOpt{}
		connOpt.SetDefault()
	}

	err = info.Verify()
	if nil != err {
		return
	}

	err = connOpt.Verify()
	if nil != err {
		return
	}

	var db *sql.DB
	db, err = sql.Open(info.Driver, info.MakeDBSource())
	if nil != err {
		return
	}

	db.SetMaxOpenConns(int(connOpt.MaxConnOpen))
	db.SetMaxIdleConns(int(connOpt.MaxConnIdle))

	err = db.Ping()
	if nil != err {
		db.Close()
		return
	}

	klog.LogInfo("[database][driver:%v][host:%v][port:%v][db:%v][maxconn:%v][maxidle:%v] Database connection open succeed",
		info.Driver, info.Host, info.Port, info.Database, connOpt.MaxConnOpen, connOpt.MaxConnIdle)

	object = &KDB{
		KObject:		kobject.NewKObject("KDB"),
		db:				db,
		info:			info.Clone(),
		connOpt:		connOpt.Clone(),
	}

	go object.reporting()

	return
}

func (m *KDB) Destroy() {

	if err := m.Close() ; nil != err {
		klog.LogWarn("KDB.Destroy() db close err : %s", err.Error())
	}

	m.KObject.Destroy()
}

func (m *KDB) Close() (err error) {

	err = m.db.Close()
	return
}

func (m *KDB) Exec(query string, args ...interface{}) (result sql.Result, err error) {

	timer := kutil.KTimer{}
	if m.connOpt.ResponseTimeCheck {
		timer.Reset()
	}

	result, err = m.db.Exec(query, args...)

	if m.connOpt.ResponseTimeCheck {
		elapsed := uint32(timer.ElapsedMilisec())
		if elapsed > m.connOpt.ResponseTimeLimit {
			klog.LogWarn("[database] Response time delayed, query:%s, time:%d", query, elapsed)
		}
	}

	return
}

func (m *KDB) QueryRow(query string, args ...interface{}) (row *sql.Row) {

	timer := kutil.KTimer{}
	if m.connOpt.ResponseTimeCheck {
		timer.Reset()
	}

	row = m.db.QueryRow(query, args...)

	if m.connOpt.ResponseTimeCheck {
		elapsed := uint32(timer.ElapsedMilisec())
		if elapsed > m.connOpt.ResponseTimeLimit {
			klog.LogWarn("[database] Response time delayed, query:%s, time:%d", query, elapsed)
		}
	}

	return
}

func (m *KDB) Query(query string, args ...interface{}) (rows *sql.Rows) {

	var err error
	timer := kutil.KTimer{}
	if m.connOpt.ResponseTimeCheck {
		timer.Reset()
	}

	rows, err = m.db.Query(query, args...)

	if m.connOpt.ResponseTimeCheck {
		elapsed := uint32(timer.ElapsedMilisec())
		if elapsed > m.connOpt.ResponseTimeLimit {
			klog.LogWarn("[database] Response time delayed, query:%s, time:%d", query, elapsed)
		}
	}

	if nil != err {
		if nil != rows {
			rows.Close()
		}
		rows = nil
		klog.LogWarn("[database] Query error, query:%s, err:%s, args:%v", query, err.Error(), args)
	}
	return
}


func (m *KDB) reporting() {

	defer func() {
		if rc := recover() ; nil != rc {
			klog.LogFatal("KDB.reporting() recovered : %v", rc)
		}
	}()

	interval := time.Duration(m.connOpt.ReportingInterval)*time.Millisecond

	if 0 >= interval {
		return
	}

	timer := time.NewTimer(interval)

	for {

		select {
		case <-m.DestroySignal():
			klog.LogDetail("KDB.reporting() Destroy sensed")
			return
		case <-timer.C:
			klog.LogInfo("[host:%s][db:%s][openconnection:%d]", m.info.Host, m.info.Database, m.db.Stats().OpenConnections)
			timer.Reset(interval)
		}

	}
}




package kdatabase

import (

	"database/sql"
	"kobject"
	"errors"

	klog "klogger"
	"kutil"
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

	object = &KDB{
		KObject:		kobject.NewKObject("KDB"),
		db:				db,
		info:			info.Clone(),
		connOpt:		connOpt.Clone(),
	}

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
		rows = nil
		klog.LogWarn("[database] Query error, query:%s, err:%s, args:%v", query, err.Error(), args)
	}
	return
}

func (m *KDB) QueryResults(query string, args ...interface{}) (result *KDBResult) {

	var rows *sql.Rows
	rows = m.Query(query, args...)
	if nil == rows {
		return
	}
	defer rows.Close()

	result = NewKDBResult()

	for {
		colnames, err := rows.Columns()
		if nil != err {
			result = nil
			klog.LogWarn("[database] QueryResults rows.Columns error, query:%s, err:%s, args:%v", query, err.Error(), args)
			break
		}

		set := NewKDBSet()
		colcount	:= len(colnames)

		for rows.Next() {
			records 	:= make([]interface{},colcount)
			err = rows.Scan(records...)
			if nil != err {
				result = nil
				klog.LogWarn("[database] QueryResults rows.Scan error, query:%s, err:%s, args:%v", query, err.Error(), args)
				return
			}
			set.rows = append(set.rows, records)
		}

		result.sets = append(result.sets, set)

		if false == rows.NextResultSet() {
			break
		}
	}

	return
}




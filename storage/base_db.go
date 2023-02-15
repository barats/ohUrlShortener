// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package storage

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"ohurlshortener/utils"
)

var dbService = &DatabaseService{}

// DatabaseService 数据库服务
type DatabaseService struct {
	Connection *sqlx.DB
}

// InitDatabaseService 初始化数据库服务
func InitDatabaseService() (*DatabaseService, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		utils.DatabaseConfig.Host, utils.DatabaseConfig.Port, utils.DatabaseConfig.User,
		utils.DatabaseConfig.Password, utils.DatabaseConfig.DbName)
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return dbService, err
	}
	conn.SetMaxOpenConns(utils.DatabaseConfig.MaxOpenConns)
	conn.SetMaxIdleConns(utils.DatabaseConfig.MaxIdleConn)
	conn.SetConnMaxLifetime(0) // always REUSE
	dbService.Connection = conn
	return dbService, nil
}

// DbNamedExec 执行带有命名参数的sql语句
func DbNamedExec(query string, args interface{}) error {
	_, err := dbService.Connection.NamedExec(query, args)
	return err
}

// DbExecTx 执行事务
func DbExecTx(query ...string) error {
	tx := dbService.Connection.MustBegin()
	for _, s := range query {
		tx.MustExec(s)
	} // end of for
	err := tx.Commit()
	if err != nil {
		return tx.Rollback()
	}
	return nil
} // end of func

//
// func DbExecTx(query string, args ...interface{}) error {
//	tx, err := dbService.Connection.Begin()
//	if err != nil {
//		return err
//	}
//	defer tx.Commit()
//
//	stmt, err := tx.Prepare(dbService.Connection.Rebind(query))
//	if err != nil {
//		return err
//	}
//	defer stmt.Close()
//
//	_, err = stmt.Exec(args...)
//	if err != nil {
//		return err
//	}
//
//	return nil
// }

// DbGet 获取单条记录
func DbGet(query string, dest interface{}, args ...interface{}) error {
	err := dbService.Connection.Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// DbSelect 获取多条记录
func DbSelect(query string, dest interface{}, args ...interface{}) error {
	return dbService.Connection.Select(dest, query, args...)
}

// DbClose 关闭数据库连接
func DbClose() {
	dbService.Connection.Close()
}

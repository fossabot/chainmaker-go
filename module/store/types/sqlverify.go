/*
 * Copyright (C) BABEC. All rights reserved.
 * Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package types

import (
	"errors"
	"regexp"
	"strings"
)

var ERROR_INVALID_SQL = errors.New("invalid sql")
var ERROR_FORBIDDEN_SQL = errors.New("forbidden sql")

//如果状态数据库是标准SQL语句，对标准SQL的SQL语句进行语法检查，不关心具体的SQL DB类型的语法差异
type StandardSqlVerify struct {
}

func (s *StandardSqlVerify) VerifyDDLSql(sql string) error {
	SQL := strings.ToUpper(sql)
	if err := s.checkForbiddenSql(SQL); err != nil {
		return err
	}
	reg := regexp.MustCompile(`^(CREATE|ALTER|DROP)\s+(TABLE|VIEW|INDEX)`)
	match := reg.MatchString(SQL)
	if match {
		return nil
	}
	if strings.HasPrefix(SQL, "TRUNCATE TABLE") {
		return nil
	}
	return ERROR_INVALID_SQL

}
func (s *StandardSqlVerify) VerifyDMLSql(sql string) error {
	SQL := strings.ToUpper(sql)
	if err := s.checkForbiddenSql(SQL); err != nil {
		return err
	}
	reg := regexp.MustCompile(`^(INSERT|UPDATE|DELETE)\s+`)
	match := reg.MatchString(SQL)
	if match {
		return nil
	}
	return ERROR_INVALID_SQL
}
func (s *StandardSqlVerify) VerifyDQLSql(sql string) error {
	SQL := strings.ToUpper(sql)
	if err := s.checkForbiddenSql(SQL); err != nil {
		return err
	}
	reg := regexp.MustCompile(`^SELECT\s+`)
	match := reg.MatchString(SQL)
	if match {
		return nil
	}
	return ERROR_INVALID_SQL
}

//禁用use database,禁用 select * from anotherdb.table形式
func (s *StandardSqlVerify) checkForbiddenSql(sql string) error {
	SQL := strings.ToUpper(sql)
	reg := regexp.MustCompile(`^USE\s+`)
	match := reg.MatchString(SQL)
	if match {
		return ERROR_FORBIDDEN_SQL
	}
	tableNames := s.getSqlTableName(sql)
	for _, tableName := range tableNames {
		if strings.Contains(tableName, ".") {
			return ERROR_FORBIDDEN_SQL
		}
	}
	return nil
}

//用于测试场景，不对SQL语句进行检查，任意SQL检查都通过
type SqlVerifyPass struct {
}

func (s *SqlVerifyPass) VerifyDDLSql(sql string) error {
	return nil
}
func (s *SqlVerifyPass) VerifyDMLSql(sql string) error {
	return nil
}
func (s *SqlVerifyPass) VerifyDQLSql(sql string) error {
	return nil
}

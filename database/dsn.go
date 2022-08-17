package database

import (
	"regexp"

	"github.com/zedisdog/sweetbean/errx"
	"github.com/zedisdog/sweetbean/tools"
)

// Type is database type. eg. mysql, postgres
type Type string

const (
	TypeMysql    Type = "mysql"
	TypePostgres Type = "postgres"
)

type DSN string

func (d DSN) Encode() string {
	return tools.EncodeQuery(string(d))
}

func (d DSN) split() []string {
	reg := regexp.MustCompile(`(^\S+)://(\S+$)`)
	info := reg.FindStringSubmatch(d.Encode())
	if len(info) < 3 {
		panic(errx.New("dsn is invalid, forget schema?"))
	}
	return info[1:]
}

func (d DSN) Type() Type {
	return Type(d.split()[0])
}

func (d DSN) RemoveSchema() string {
	return d.split()[1]
}

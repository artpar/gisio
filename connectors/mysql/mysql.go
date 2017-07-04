package mysql

import (
  "database/sql"
  "github.com/ziutek/mymysql/mysql"
  "fmt"
)

type ConnectionConfig struct {
  Hostname string
  Port     string
  Username string
  Password string
  DbName   string
}

func NewConnection(c ConnectionConfig) (*sql.DB) {
  db := mysql.New("tcp", "", fmt.Sprintf("%s:%s", c.Hostname, c.Port), c.Username, c.Password, c.DbName)
  return db
}

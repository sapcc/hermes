package configdb

import (
	"database/sql"
	"github.com/sapcc/hermes/pkg/util"
	"github.com/spf13/viper"
	//enable mysql driver for database/sql
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"
)

//DB holds the main database connection. It will be `nil` until InitDatabase() is called.
type MySQL struct {
	DB *gorp.DbMap
}

func (mysql *MySQL) dbmap() *gorp.DbMap {
	//Lazy initialization
	if mysql.DB == nil {
		mysql.init()
	}
	return mysql.DB
}

func (mysql *MySQL) init() {

	var err error
	var dsn = viper.GetString("mysql.dsn")
	util.LogDebug("Using MySQL DSN: %s", dsn)
	sqlDriver := "mysql"
	db, err := sql.Open(sqlDriver, dsn)
	if err != nil {
		panic(err)
	}
	mysql.DB = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
}

/* cases for action
1. Get Audit Config -> Doesn't exist in DB
	Return False
2. Get Audit Config -> Does exist
	Return Value
*/
func (mysql MySQL) GetAudit(tenantId string) (*AuditConfig, error) {
	dbmap := mysql.dbmap()

	var auditdetail AuditConfig
	err := dbmap.SelectOne(&auditdetail, "select enabled, tenant_id from auditconfig where tenant_id=?", tenantId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Didn't find tenant_id in the database
			auditdetail.Enabled = false
			auditdetail.TenantID = tenantId
		} else {
			util.LogError("AuditConfig query failed: %v", err)
			return nil, err
		}
	}

	return &auditdetail, nil
}

/*
	1. Put Audit Config -> Doesn't exist in DB
		insert new tenant_id/enabled
	2. Put Audit Config -> Does exist in DB
		update with new value
*/
func (mysql MySQL) PutAudit(tenantId string) (*AuditConfig, error) {
	dbmap := mysql.dbmap()
	if tenantId == "" {
		util.LogError("No TenantId in Put Request")
		//return nil, nil
	}

	var auditdetail AuditConfig
	dbmap.AddTableWithName(auditdetail, "auditconfig")
	err := dbmap.SelectOne(&auditdetail, "select enabled, tenant_id from auditconfig where tenant_id=?", tenantId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Didn't find tenant_id in the database, need to create it.
			auditdetail.Enabled = false
			auditdetail.TenantID = tenantId
			err := dbmap.Insert(&auditdetail)
			if err != nil {
				util.LogError("Problem inserting row to MySQL: %v", err)
				return nil, err
			}
		} else {
			util.LogError("AuditConfig query failed: %v", err)
			return nil, err
		}
	}

	return &auditdetail, nil
}

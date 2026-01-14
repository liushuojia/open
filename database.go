package utils

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type (
	DatabaseTable struct {
		createSqlString string
		FieldList       []string
		FieldMap        map[string]*DatabaseTableField
	}
	DatabaseTableField struct {
		Field      string `gorm:"column:Field"`
		Type       string `gorm:"column:Type"`
		Collation  string `gorm:"column:Collation"`
		Null       string `gorm:"column:Null"`
		Key        string `gorm:"column:Key"`
		Default    string `gorm:"column:Default"`
		Extra      string `gorm:"column:Extra"`
		Privileges string `gorm:"column:Privileges"`
		Comment    string `gorm:"column:Comment"`

		PrimaryKey    bool   `gorm:"-"`
		Unique        bool   `gorm:"-"`
		NotNull       bool   `gorm:"-"`
		AutoIncrement bool   `gorm:"-"`
		Index         bool   `gorm:"-"`
		Unsigned      bool   `gorm:"-"`
		TypeLength    string `gorm:"-"`
		TypeString    string `gorm:"-"`
		GOType        string `gorm:"-"`
		ProtoType     string `gorm:"-"`
	}
)

type Database struct {
	db *gorm.DB

	TableList []string
	TableMap  map[string]*DatabaseTable
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) SetDB(db *gorm.DB) *Database {
	d.db = db
	return d
}
func (d *Database) Get() error {
	rows, err := d.db.Raw("show tables").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	d.TableList = make([]string, 0)
	d.TableMap = make(map[string]*DatabaseTable)

	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			continue
		}

		for i := range cols {
			val := columnPointers[i].(*interface{})

			tableName := ""
			if fmt.Sprintf("%T", *val) == "[]uint8" {
				tableName = string((*val).([]uint8))
			}
			if fmt.Sprintf("%T", *val) == "int64" {
				tableName = strconv.Itoa(int((*val).(int64)))
			}
			if tableName == "" {
				continue
			}
			d.TableList = append(d.TableList, tableName)

			databaseTable, err := d.GetTable(tableName)
			if err != nil {
				return err
			}
			d.TableMap[tableName] = databaseTable
		}
	}

	return nil
}
func (d *Database) GetTable(tableName string) (*DatabaseTable, error) {
	databaseTable := &DatabaseTable{}

	rows, err := d.db.Raw(fmt.Sprintf("SHOW CREATE table `%s`", tableName)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	for rows.Next() {
		var (
			columns        = make([]interface{}, len(cols))
			columnPointers = make([]interface{}, len(cols))
		)

		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			continue
		}

		for i, colName := range cols {
			var (
				val   = columnPointers[i].(*interface{})
				value = ""
			)

			if fmt.Sprintf("%T", *val) == "[]uint8" {
				value = string((*val).([]uint8))
			}
			if fmt.Sprintf("%T", *val) == "int64" {
				value = strconv.Itoa(int((*val).(int64)))
			}

			if value == "" {
				continue
			}
			switch colName {
			case "Create Table":
				databaseTable.createSqlString = value
			}
		}
	}

	databaseTable.FieldList, databaseTable.FieldMap, err = d.GetTableField(tableName)
	if err != nil {
		return nil, err
	}

	return databaseTable, nil
}
func (d *Database) GetTableField(tableName string) ([]string, map[string]*DatabaseTableField, error) {
	var (
		fieldList    = make([]string, 0)
		fieldMap     = make(map[string]*DatabaseTableField)
		fieldSQLList = make([]*DatabaseTableField, 0)
	)

	if err := d.db.Raw(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", tableName)).Scan(&fieldSQLList).Error; err != nil {
		return nil, nil, err
	}

	for _, vv := range fieldSQLList {
		vv.AutoIncrement = vv.Extra == "auto_increment"
		vv.NotNull = vv.Null == "NO"

		switch vv.Key {
		case "PRI":
			vv.PrimaryKey = true
		case "UNI":
			vv.Unique = true
		case "MUL":
			vv.Index = true
		}

		var (
			typeString = strings.ToLower(vv.Type)
			typeLength = ""
		)

		tmp1 := strings.Split(typeString, "(")
		if len(tmp1) > 1 {
			typeString = tmp1[0]
			typeLength = strings.Trim(tmp1[1], ")")
		}
		vv.TypeString = typeString
		vv.TypeLength = typeLength

		if strings.Index(vv.Type, "unsigned") != -1 {
			vv.Unsigned = true
		}

		vv.ProtoType = vv.TypeString
		switch strings.ToLower(vv.TypeString) {
		case "timestamp":
			vv.ProtoType = "uint64"
		case "decimal":
			vv.ProtoType = "decimal"
		case "bigint":
			vv.ProtoType = "int64"
			if vv.Unsigned {
				vv.ProtoType = "uint64"
			}
		case "int":
			vv.ProtoType = "int32"
			if vv.Unsigned {
				vv.ProtoType = "uint32"
			}
		case "char":
			vv.ProtoType = "string"
		case "tinyint":
			if vv.TypeLength == "1" {
				vv.ProtoType = "bool"
				break
			}
			vv.ProtoType = "int32"
			if vv.Unsigned {
				vv.ProtoType = "uint32"
			}
		case "text":
			str := "_id_list"
			if vv.Field[len(vv.Field)-len(str):] == str {
				vv.ProtoType = "repeated uint64"
				break
			}
			str = "_list"
			if vv.Field[len(vv.Field)-len(str):] == str {
				vv.ProtoType = "repeated string"
				break
			}
			vv.ProtoType = "string"
		}

		fieldMap[vv.Field] = vv
		fieldList = append(fieldList, vv.Field)
	}

	return fieldList, fieldMap, nil
}

// ResetTableName a-bb_cc -> ABbCc
func (d *Database) ResetTableName(name string) string {
	if len(name) <= 2 {
		return strings.ToUpper(name)
	}
	tmp := strings.Split(strings.Replace(name, "-", "_", -1), "_")
	for k, v := range tmp {
		tmp[k] = strings.ToUpper(v[:1]) + v[1:]
	}
	return strings.Join(tmp, "")
}

// ResetFieldName a-bb_cc -> aBbCc
func (d *Database) ResetFieldName(name string) string {
	if len(name) <= 2 {
		return strings.ToUpper(name)
	}
	tmp := strings.Split(strings.Replace(name, "-", "_", -1), "_")
	for k, v := range tmp {
		if k == 0 {
			continue
		}
		tmp[k] = strings.ToUpper(v[:1]) + v[1:]
	}
	return strings.Join(tmp, "")
}

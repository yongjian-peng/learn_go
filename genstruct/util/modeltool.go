package util

import (
	"html/template"
	"strings"
)

type ModelInfo struct {
	BDName       string
	DBConnection string
	TableName    string
	PackageName  string
	ModelName    string
	TableSchema  *[]TABLE_SCHEMA
}

type TABLE_SCHEMA struct {
	COLUMN_NAME    string `db:"COLUMN_NAME" json:"column_name"`
	DATA_TYPE      string `db:"DATA_TYPE" json:"data_type"`
	COLUMN_KEY     string `db:"COLUMN_KEY" json:"column_key"`
	COLUMN_COMMENT string `db:"COLUMN_COMMENT" json:"column_comment"`
	IS_NULLABLE    string `db:"IS_NULLABLE" json:"is_nullable"`
	COLUMN_TYPE    string `db:"COLUMN_TYPE" json:"column_type"`
	EXTRA          string `db:"EXTRA" json:"extra"`
}

func (m *ModelInfo) ColumnNames() []string {
	result := make([]string, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {

		result = append(result, t.COLUMN_NAME)

	}
	return result
}

func (m *ModelInfo) ColumnCount() int {
	return len(*m.TableSchema)
}

func (m *ModelInfo) PkColumnsSchema() []TABLE_SCHEMA {
	result := make([]TABLE_SCHEMA, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		if t.COLUMN_KEY == "PRI" {
			result = append(result, t)
		}
	}
	return result
}

func (m *ModelInfo) HavePk() bool {
	return len(m.PkColumnsSchema()) > 0
}

func (m *ModelInfo) NoPkColumnsSchema() []TABLE_SCHEMA {
	result := make([]TABLE_SCHEMA, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		if t.COLUMN_KEY != "PRI" {
			result = append(result, t)
		}
	}
	return result
}

func (m *ModelInfo) NoPkColumns() []string {
	noPkColumnsSchema := m.NoPkColumnsSchema()
	result := make([]string, 0, len(noPkColumnsSchema))
	for _, t := range noPkColumnsSchema {
		result = append(result, t.COLUMN_NAME)
	}
	return result
}

func (m *ModelInfo) PkColumns() []string {
	pkColumnsSchema := m.PkColumnsSchema()
	result := make([]string, 0, len(pkColumnsSchema))
	for _, t := range pkColumnsSchema {
		result = append(result, t.COLUMN_NAME)
	}
	return result
}

func IsUUID(str string) bool {
	return "uuid" == str
}

func FirstCharLower(str string) string {
	if len(str) > 0 {
		return strings.ToLower(str[0:1]) + str[1:]
	} else {
		return ""
	}
}

func FirstCharUpper(str string) string {
	showName := ""
	if strings.Contains(str, "_") {
		names := strings.Split(str, "_")
		for _, nameItem := range names {
			showName += strFirstToUpper(nameItem)
		}
	} else {
		showName = strFirstToUpper(str)
	}
	return showName
}

func toCamelCase(s string) string {
	var result string
	var nextUpper bool

	for _, c := range s {
		if c == '_' || c == '-' {
			nextUpper = true
			continue
		}
		if nextUpper {
			result += strings.ToUpper(string(c))
			nextUpper = false
		} else {
			result += string(c)
		}
	}

	return result
}

func SecondCharUpper(str string) string {
	showName := ""
	if strings.Contains(str, "_") {
		names := strings.Split(str, "_")
		for index, nameItem := range names {
			if index == 0 {
				showName += nameItem
			} else {
				showName += strFirstToUpper(nameItem)
			}
		}
	} else {
		showName = str
	}
	return showName
}

/**
 * 字符串首字母转化为大写 ios_bbbbbbbb -> iosBbbbbbbbb
 */
func strFirstToUpper(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			vv[i] -= 32
			upperStr += string(vv[i]) // + string(vv[i+1])
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func Tags(columnName string, isNullable string, columnType string, columnKey string, extra string) template.HTML {

	tag := "`gorm:" + `"column:` + columnName + ";type:" + columnType
	if columnKey == "PRI" {
		tag += ";primary_key"
	}
	if extra != "" {
		tag += ";" + strings.ToUpper(extra)
	}
	if isNullable == "NO" {
		tag += ";NOT NULL"
	}
	tag += `"` + " json:" + `"` + toCamelCase(columnName) + "\"`"

	return template.HTML(tag)
}

func ExportColumn(columnName string) string {
	columnItems := strings.Split(columnName, "_")
	columnItems[0] = FirstCharUpper(columnItems[0])
	for i := 0; i < len(columnItems); i++ {
		item := strings.Title(columnItems[i])

		if strings.ToUpper(item) == "ID" {
			item = "Id"
		}

		columnItems[i] = item
	}

	return strings.Join(columnItems, "")

}

func TypeConvert(str string, isnull string) string {
	switch str {
	case "smallint", "tinyint":
		return "int32"
	case "varchar", "text", "longtext", "char":
		return "string"
	case "date":
		if isnull == "YES" {
			return "sql.NullTime"
		}
		return "time.Time"
	case "int":
		return "int32"
	case "timestamp", "datetime":
		if isnull == "YES" {
			return "sql.NullTime"
		}
		return "time.Time"
	case "bigint":
		return "int64"
	case "float", "double", "decimal":
		return "float64"
	default:
		return str
	}
}

func Join(a []string, sep string) string {
	return strings.Join(a, sep)
}

func ColumnAndType(table_schema []TABLE_SCHEMA) string {
	result := make([]string, 0, len(table_schema))
	for _, t := range table_schema {
		result = append(result, t.COLUMN_NAME+" "+TypeConvert(t.DATA_TYPE, t.IS_NULLABLE))
	}
	return strings.Join(result, ",")
}

func ColumnWithPostfix(columns []string, Postfix, sep string) string {
	result := make([]string, 0, len(columns))
	for _, t := range columns {
		result = append(result, t+Postfix)
	}
	return strings.Join(result, sep)
}

func MakeQuestionMarkList(num int) string {
	a := strings.Repeat("?,", num)
	return a[:len(a)-1]
}

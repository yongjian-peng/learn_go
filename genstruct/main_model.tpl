{{$exportModelName := .ModelName | FirstCharUpper}}

package {{.PackageName}}

type {{$exportModelName}} struct {
{{range .TableSchema}} {{.COLUMN_NAME | ExportColumn}} {{TypeConvert .DATA_TYPE  .IS_NULLABLE}} {{Tags .COLUMN_NAME .IS_NULLABLE .COLUMN_TYPE .COLUMN_KEY .EXTRA }} // {{.COLUMN_COMMENT}}
{{end}}}


func (m *{{$exportModelName}}) TableName() string {
	return "{{.TableName}}"
}
package goutils

const (
	TimeLayout   = "2006-01-02 15:04:05"
	TimeLayout_2 = "20060102150405"
	DateLayout   = "2006-01-02"
	NULL         = ""
)

type File struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

/**
 * CopyMap 赋值map
 */
func CopyMap(m map[string]interface{}) map[string]interface{} {
	m2 := make(map[string]interface{}, len(m))
	for k, v := range m {
		m2[k] = v
	}
	// id should not be accessible here, it should exist only inside loop
	return m2
}

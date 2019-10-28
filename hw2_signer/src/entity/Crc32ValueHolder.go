package entity

type Crc32ValueHolder struct {
	value      string
	crc32Value string
}

func (crc32ValueHolder *Crc32ValueHolder) GetCrc32Value() string {
	return crc32ValueHolder.crc32Value
}

func (crc32ValueHolder *Crc32ValueHolder) GetValue() string {
	return crc32ValueHolder.value
}

func (crc32ValueHolder *Crc32ValueHolder) SetCrc32Value(crc32Value string) {
	crc32ValueHolder.crc32Value = crc32Value
}

func (crc32ValueHolder *Crc32ValueHolder) SetValue(value string) {
	crc32ValueHolder.value = value
}

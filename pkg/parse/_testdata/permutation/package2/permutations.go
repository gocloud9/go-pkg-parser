package package2

type SomeStruct struct {
	StringField         string               `json:"string_field" yaml:"stringField"`
	IntField            int                  `json:"int_field" yaml:"intField"`
	BoolField           bool                 `json:"bool_field" yaml:"boolField"`
	ChanField           chan int             `json:"chan_field" yaml:"chanField"`
	MapField            map[string]int       `json:"map_field" yaml:"mapField"`
	SliceField          []int                `json:"slice_field" yaml:"sliceField"`
	SubStructField      SubStruct            `json:"sub_struct_field" yaml:"subStructField"`
	SubStructMapField   map[string]SubStruct `json:"sub_struct_map_field" yaml:"subStructMapField"`
	SubStructSliceField []SubStruct          `json:"sub_struct_slice_field" yaml:"subStructSliceField"`
}

type SubStruct struct {
	Something string
}

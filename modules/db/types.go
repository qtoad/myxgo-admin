// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"html/template"
	"strconv"
)

// FieldType is the database field type.
type FieldType string

const (
	// =================================
	// integer
	// =================================

	Int       FieldType = "INT"
	Tinyint   FieldType = "TINYINT"
	Mediumint FieldType = "MEDIUMINT"
	Smallint  FieldType = "SMALLINT"
	Bigint    FieldType = "BIGINT"
	Bit       FieldType = "BIT"
	Int8      FieldType = "INT8"
	Int4      FieldType = "INT4"
	Int2      FieldType = "INT2"

	Integer     FieldType = "INTEGER"
	Numeric     FieldType = "NUMERIC"
	Smallserial FieldType = "SMALLSERIAL"
	Serial      FieldType = "SERIAL"
	Bigserial   FieldType = "BIGSERIAL"
	Money       FieldType = "MONEY"

	// =================================
	// float
	// =================================

	Real    FieldType = "REAL"
	Float   FieldType = "FLOAT"
	Float4  FieldType = "FLOAT4"
	Float8  FieldType = "FLOAT8"
	Double  FieldType = "DOUBLE"
	Decimal FieldType = "DECIMAL"

	Doubleprecision FieldType = "DOUBLEPRECISION"

	// =================================
	// string
	// =================================

	Date      FieldType = "DATE"
	Time      FieldType = "TIME"
	Year      FieldType = "YEAR"
	Datetime  FieldType = "DATETIME"
	Timestamp FieldType = "TIMESTAMP"

	Text       FieldType = "TEXT"
	Longtext   FieldType = "LONGTEXT"
	Mediumtext FieldType = "MEDIUMTEXT"
	Tinytext   FieldType = "TINYTEXT"

	Varchar FieldType = "VARCHAR"
	Char    FieldType = "CHAR"
	Bpchar  FieldType = "BPCHAR"
	JSON    FieldType = "JSON"

	Blob       FieldType = "BLOB"
	Tinyblob   FieldType = "TINYBLOB"
	Mediumblob FieldType = "MEDIUMBLOB"
	Longblob   FieldType = "LONGBLOB"

	Interval FieldType = "INTERVAL"
	Boolean  FieldType = "BOOLEAN"
	Bool     FieldType = "BOOL"

	Point   FieldType = "POINT"
	Line    FieldType = "LINE"
	Lseg    FieldType = "LSEG"
	Box     FieldType = "BOX"
	Path    FieldType = "PATH"
	Polygon FieldType = "POLYGON"
	Circle  FieldType = "CIRCLE"

	Cidr    FieldType = "CIDR"
	Inet    FieldType = "INET"
	Macaddr FieldType = "MACADDR"

	Character        FieldType = "CHARACTER"
	Varyingcharacter FieldType = "VARYINGCHARACTER"
	Nchar            FieldType = "NCHAR"
	Nativecharacter  FieldType = "NATIVECHARACTER"
	Nvarchar         FieldType = "NVARCHAR"
	Clob             FieldType = "CLOB"

	Binary    FieldType = "BINARY"
	Varbinary FieldType = "VARBINARY"
	Enum      FieldType = "ENUM"
	Set       FieldType = "SET"

	Geometry FieldType = "GEOMETRY"

	Multilinestring    FieldType = "MULTILINESTRING"
	Multipolygon       FieldType = "MULTIPOLYGON"
	Linestring         FieldType = "LINESTRING"
	Multipoint         FieldType = "MULTIPOINT"
	Geometrycollection FieldType = "GEOMETRYCOLLECTION"

	Name FieldType = "NAME"
	UUID FieldType = "UUID"

	Timestamptz FieldType = "TIMESTAMPTZ"
	Timetz      FieldType = "TIMETZ"
)

// GetFieldType turn the string value into FieldType.
func GetFieldType(s string) FieldType {
	return FieldType(s)
}

// GetFieldTypeAndCheck check the FieldType.
func GetFieldTypeAndCheck(s string) FieldType {
	ss := FieldType(s)
	if !Contains(ss, BoolTypeList) &&
		!Contains(ss, IntTypeList) &&
		!Contains(ss, FloatTypeList) &&
		!Contains(ss, UintTypeList) &&
		!Contains(ss, StringTypeList) {
		panic("wrong type: " + s)
	}
	return ss
}

var (
	// StringTypeList is a FieldType list of string.
	StringTypeList = []FieldType{Date, Time, Year, Datetime, Timestamptz, Timestamp, Timetz,
		Varchar, Char, Mediumtext, Longtext, Tinytext,
		Text, JSON, Blob, Tinyblob, Mediumblob, Longblob,
		Interval, Point, Bpchar,
		Line, Lseg, Box, Path, Polygon, Circle, Cidr, Inet, Macaddr, Character, Varyingcharacter,
		Nchar, Nativecharacter, Nvarchar, Clob, Binary, Varbinary, Enum, Set, Geometry, Multilinestring,
		Multipolygon, Linestring, Multipoint, Geometrycollection, Name, UUID, Timestamptz,
		Name, UUID, Inet}

	// BoolTypeList is a FieldType list of bool.
	BoolTypeList = []FieldType{Bool, Boolean}

	// IntTypeList is a FieldType list of integer.
	IntTypeList = []FieldType{Int4, Int2, Int8,
		Int,
		Tinyint,
		Mediumint,
		Smallint,
		Smallserial, Serial, Bigserial,
		Integer,
		Bigint}

	// FloatTypeList is a FieldType list of float.
	FloatTypeList = []FieldType{Float, Float4, Float8, Double, Real, Doubleprecision}

	// UintTypeList is a FieldType list of uint.
	UintTypeList = []FieldType{Decimal, Bit, Money, Numeric}
)

// Contains check the given FieldType is in the list or not.
func Contains(v FieldType, a []FieldType) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}

// Value is a string.
type Value string

// ToInt64 turn the string to a int64.
func (v Value) ToInt64() int64 {
	value, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		panic("wrong value")
	}
	return value
}

// String return the string value.
func (v Value) String() string {
	return string(v)
}

// HTML return the template.HTML value.
func (v Value) HTML() template.HTML {
	return template.HTML(v)
}

// GetValueFromFieldType return Value of given FieldType and interface.
func GetValueFromFieldType(typ FieldType, value interface{}, json bool) Value {
	if json {
		return GetValueFromJSONFieldType(typ, value)
	} else {
		return GetValueFromSQLFieldType(typ, value)
	}
}

// GetValueFromSQLFieldType return Value of given FieldType and interface.
func GetValueFromSQLFieldType(typ FieldType, value interface{}) Value {
	switch {
	case Contains(typ, StringTypeList):
		if v, ok := value.(string); ok {
			return Value(v)
		}
		return ""
	case Contains(typ, BoolTypeList):
		if v, ok := value.(bool); ok {
			if v {
				return "true"
			}
			return "false"
		}
		if v, ok := value.(int64); ok {
			if v == 0 {
				return "false"
			}
			return "true"
		}
		return "false"
	case Contains(typ, IntTypeList):
		if v, ok := value.(int64); ok {
			return Value(fmt.Sprintf("%d", v))
		}
		return "0"
	case Contains(typ, FloatTypeList):
		if v, ok := value.(float64); ok {
			return Value(fmt.Sprintf("%f", v))
		}
		return "0"
	case Contains(typ, UintTypeList):
		if v, ok := value.([]uint8); ok {
			return Value(string(v))
		}
		return "0"
	}
	panic("wrong type：" + string(typ))
}

// GetValueFromJSONFieldType return Value of given FieldType and interface from JSON string value.
func GetValueFromJSONFieldType(typ FieldType, value interface{}) Value {
	switch {
	case Contains(typ, StringTypeList):
		if v, ok := value.(string); ok {
			return Value(v)
		}
		return ""
	case Contains(typ, BoolTypeList):
		if v, ok := value.(bool); ok {
			if v {
				return "true"
			}
			return "false"
		}
		return "false"
	case Contains(typ, IntTypeList):
		if v, ok := value.(float64); ok {
			return Value(fmt.Sprintf("%d", int64(v)))
		}
		return Value(fmt.Sprintf("%d", value))
	case Contains(typ, FloatTypeList):
		return Value(fmt.Sprintf("%f", value))
	case Contains(typ, UintTypeList):
		if v, ok := value.([]uint8); ok {
			return Value(string(v))
		}
		return "0"
	}
	panic("wrong type：" + string(typ))
}

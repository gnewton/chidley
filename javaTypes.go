package main

const JavaString = "String"
const JavaBoolean = "boolean"
const JavaShort = "short"
const JavaFloat = "float"
const JavaDouble = "double"
const JavaInt = "int"
const JavaLong = "long"

const GoBool = "bool"

const GoInt8 = "int8"
const GoUint8 = "uint8"
const GoInt16 = "int16"

const GoUint16 = "uint16"
const GoInt32 = "int32"

const GoUint32 = "uint32"
const GoInt64 = "int64"

const GoFloat32 = "float32"

const GoFloat64 = "float64"

func findJavaType(nti *NodeTypeInfo, useType bool) string {
	if !useType {
		return JavaString
	}
	goType := findType(nti, true)

	switch goType {
	case GoBool:
		return JavaBoolean
	case GoInt8, GoUint8, GoInt16:
		return JavaShort
	case GoUint16, GoInt32:
		return JavaInt
	case GoUint32, GoInt64:
		return JavaLong
	case GoFloat32:
		return JavaFloat
	case GoFloat64:
		return JavaDouble
	}
	return JavaString
}

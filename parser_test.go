package jajson_test

import (
	"github.com/aleksandrzhukovskii/jajson"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ParserSuite struct {
	suite.Suite
}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}

func (t *ParserSuite) TestParseValues() {
	tests := []struct {
		data          []byte
		expectedType  jajson.LexemType
		expectedValue []byte
	}{
		{
			data:          []byte(`   "hello world"    `),
			expectedType:  jajson.String,
			expectedValue: []byte(`"hello world"`),
		}, {
			data:          []byte(`  true  `),
			expectedType:  jajson.Bool,
			expectedValue: []byte(`true`),
		}, {
			data:          []byte(`   false    `),
			expectedType:  jajson.Bool,
			expectedValue: []byte(`false`),
		}, {
			data:          []byte(`   123    `),
			expectedType:  jajson.Int,
			expectedValue: []byte(`123`),
		}, {
			data:          []byte(`   -123    `),
			expectedType:  jajson.Int,
			expectedValue: []byte(`-123`),
		}, {
			data:          []byte(`   123.123    `),
			expectedType:  jajson.Float,
			expectedValue: []byte(`123.123`),
		}, {
			data:          []byte(`   -123.123    `),
			expectedType:  jajson.Float,
			expectedValue: []byte(`-123.123`),
		}, {
			data:          []byte(`   {  }    `),
			expectedType:  jajson.Object,
			expectedValue: []byte(`{  }`),
		}, {
			data:          []byte(`   { "test": "hello world"  }    `),
			expectedType:  jajson.Object,
			expectedValue: []byte(`{ "test": "hello world"  }`),
		}, {
			data:          []byte(`{"test":"hello world","test1":123}`),
			expectedType:  jajson.Object,
			expectedValue: []byte(`{"test":"hello world","test1":123}`),
		}, {
			data:          []byte(` [  ]   `),
			expectedType:  jajson.Array,
			expectedValue: []byte(`[  ]`),
		}, {
			data:          []byte(`   [ "test" ]  `),
			expectedType:  jajson.Array,
			expectedValue: []byte(`[ "test" ]`),
		}, {
			data:          []byte(`["test",123,true]`),
			expectedType:  jajson.Array,
			expectedValue: []byte(`["test",123,true]`),
		},
	}
	for _, test := range tests {
		typ, value, err := jajson.GetRawValue(test.data)
		t.NoError(err)
		t.Equal(test.expectedType, typ)
		t.Equal(test.expectedValue, value)
	}
}

func (t *ParserSuite) TestParseInts() {
	str := []byte("123")
	{
		val, err := jajson.GetInt[int](str)
		t.NoError(err)
		t.Equal(123, val)
	}

	{
		val, err := jajson.GetInt[int8](str)
		t.NoError(err)
		t.Equal(int8(123), val)
	}

	{
		val, err := jajson.GetInt[int16](str)
		t.NoError(err)
		t.Equal(int16(123), val)
	}

	{
		val, err := jajson.GetInt[int32](str)
		t.NoError(err)
		t.Equal(int32(123), val)
	}

	{
		val, err := jajson.GetInt[int64](str)
		t.NoError(err)
		t.Equal(int64(123), val)
	}

	{
		val, err := jajson.GetUInt[uint](str)
		t.NoError(err)
		t.Equal(uint(123), val)
	}

	{
		val, err := jajson.GetUInt[uint8](str)
		t.NoError(err)
		t.Equal(uint8(123), val)
	}

	{
		val, err := jajson.GetUInt[uint16](str)
		t.NoError(err)
		t.Equal(uint16(123), val)
	}

	{
		val, err := jajson.GetUInt[uint32](str)
		t.NoError(err)
		t.Equal(uint32(123), val)
	}

	{
		val, err := jajson.GetUInt[uint64](str)
		t.NoError(err)
		t.Equal(uint64(123), val)
	}
}

func (t *ParserSuite) TestParseFloats() {
	str := []byte("123.123")
	{
		val, err := jajson.GetFloat[float32](str)
		t.NoError(err)
		t.Equal(float32(123.123), val)
	}
	{
		val, err := jajson.GetFloat[float64](str)
		t.NoError(err)
		t.Equal(float64(123.123), val)
	}
}

func (t *ParserSuite) TestParseBool() {
	str := []byte("true")
	{
		val, err := jajson.GetBool(str)
		t.NoError(err)
		t.Equal(true, val)
	}

	str2 := []byte("false")
	{
		val, err := jajson.GetBool(str2)
		t.NoError(err)
		t.Equal(false, val)
	}
}

func (t *ParserSuite) TestParseString() {
	str := []byte(`   "hello world"    `)
	{
		val, err := jajson.GetString(str)
		t.NoError(err)
		t.Equal("hello world", val)
	}
}

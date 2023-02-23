package kql

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	tests := []struct {
		name     string
		b        Builder
		expected string
	}{
		{"Test empty", NewStatementBuilder(""), ""},
		{"Test simple literal", NewStatementBuilder("").AddLiteral("foo"), "foo"},
		{"Test simple literal ctor", NewStatementBuilder("foo"), "foo"},
		{"Test add literal", NewStatementBuilder("foo").AddLiteral("bar"), "foobar"},
		{
			"Test add int",
			NewStatementBuilder("MyTable | where i != ").AddInt(32).AddLiteral(" ;"),
			"MyTable | where i != int(32) ;",
		},
		{
			"Test add long",
			NewStatementBuilder("MyTable | where i != ").AddLong(32).AddLiteral(" ;"),
			"MyTable | where i != long(32) ;",
		},
		{
			"Test add real",
			NewStatementBuilder("MyTable | where i != ").AddReal(32.5).AddLiteral(" ;"),
			"MyTable | where i != real(32.5) ;",
		},
		{
			"Test add bool",
			NewStatementBuilder("MyTable | where i != ").AddBool(true).AddLiteral(" ;"),
			"MyTable | where i != bool(true) ;",
		},
		{
			"Test add datetime",
			NewStatementBuilder(
				"MyTable | where i != ",
			).AddDateTime(time.Date(2019, 1, 2, 3, 4, 5, 600, time.UTC)).AddLiteral(" ;"),
			"MyTable | where i != datetime(2019-01-02T03:04:05.0000006Z) ;",
		},
		{
			"Test add duration",
			NewStatementBuilder(
				"MyTable | where i != ",
			).AddTimespan(1*time.Hour + 2*time.Minute + 3*time.Second + 4*time.Microsecond).AddLiteral(" ;"),
			"MyTable | where i != timespan(01:02:03.0004000) ;",
		},
		{
			"Test add duration with days",
			NewStatementBuilder(
				"MyTable | where i != ",
			).AddTimespan(49*time.Hour + 2*time.Minute + 3*time.Second + 4*time.Microsecond).AddLiteral(" ;"),
			"MyTable | where i != timespan(2.01:02:03.0004000) ;",
		},
		{
			"Test add dynamic",
			NewStatementBuilder(
				"MyTable | where i != ",
			).AddDynamic(`{"a": 3, "b": 5.4}`).AddLiteral(" ;"),
			`MyTable | where i != dynamic("{\"a\": 3, \"b\": 5.4}") ;`,
		},
		{
			"Test add guid",
			NewStatementBuilder(
				"MyTable | where i != ",
			).AddGUID(uuid.MustParse("12345678-1234-1234-1234-123456789012")).AddLiteral(" ;"),
			"MyTable | where i != guid(12345678-1234-1234-1234-123456789012) ;",
		},
		{
			"Test add string simple",
			NewStatementBuilder(
				"MyTable | where i != ",
			).AddString("foo").AddLiteral(" ;"),
			"MyTable | where i != foo ;",
		},
		{
			"Test add string with quote",
			NewStatementBuilder(
				"MyTable | where i != ",
			).AddString("foo\"bar").AddLiteral(" ;"),
			"MyTable | where i != foo\"bar ;",
		},
		{"Test add identifiers",
			NewStatementBuilder("").
				AddDatabase("foo_1").AddLiteral(".").
				AddTable("_bar").AddLiteral(" | where ").
				AddColumn("_baz").AddLiteral(" == ").
				AddFunction("func_").AddLiteral("() ;"),
			`database("foo_1")._bar | where _baz == func_() ;`},
		{"Test add identifiers complex",
			NewStatementBuilder("").
				AddDatabase("f\"\"o").AddLiteral(".").
				AddTable("b\\a\\r").AddLiteral(" | where ").
				AddColumn("b\na\nz").AddLiteral(" == ").
				AddFunction("f_u_n\u1234c").AddLiteral("() ;"),
			`database("f\"\"o").["b\\a\\r"] | where ["b\na\nz"] == ["f_u_n\u1234c"]() ;`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.b.Build()
			assert.Equal(t, test.expected, actual.Query())
		})
	}
}

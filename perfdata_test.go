package perfdata


import (
	"fmt"
	"math"
	"testing"
)

var testPerfdataSetOk = []struct{
	in string
	out []*Perfdata
}{
	{
		in: "users=2;3;7;0",
		out: []*Perfdata{
			&Perfdata{
				Label: "users",
				UOM: "",
				Value: 2,
				Warning: 3,
				Critical: 7,
				Min: 0,
				Max: math.NaN(),
			},
		},
	},
	{
		in: "load1=0.050;7.000;10.000;0; load5=0.040;6.000;7.000;0; load15=0.010;5.000;6.000;0;100",
		out: []*Perfdata{
			&Perfdata{
				Label: "load1",
				UOM: "",
				Value: 0.050,
				Warning: 7.000,
				Critical: 10.000,
				Min: 0,
				Max: math.NaN(),
			},
			&Perfdata{
				Label: "load5",
				UOM: "",
				Value: 0.040,
				Warning: 6.000,
				Critical: 7.000,
				Min: 0,
				Max: math.NaN(),
			},
			&Perfdata{
				Label: "load15",
				UOM: "",
				Value: 0.010,
				Warning: 5.000,
				Critical: 6.000,
				Min: 0,
				Max: 100,
			},
		},
	},
	{
		in: "'users'=2%;3;7;0",
		out: []*Perfdata{
			&Perfdata{
				Label: "'users'",
				UOM: "%",
				Value: 2,
				Warning: 3,
				Critical: 7,
				Min: 0,
				Max: math.NaN(),
			},
		},
	},
	{
		in: "users=2;;;;",
		out: []*Perfdata{
			&Perfdata{
				Label: "users",
				UOM: "",
				Value: 2,
				Warning: math.NaN(),
				Critical: math.NaN(),
				Min: math.NaN(),
				Max: math.NaN(),
			},
		},
	},
	{
		in: "users=2,34",
		out: []*Perfdata{
			&Perfdata{
				Label: "users",
				UOM: "",
				Value: 2.34,
				Warning: math.NaN(),
				Critical: math.NaN(),
				Min: math.NaN(),
				Max: math.NaN(),
			},
		},
	},
}

var testPerfdataSetFail = []struct{
	in string
	err string
	out []*Perfdata
}{
	{
		in: "users'=2%;3;7;0",
		out: []*Perfdata{
			&Perfdata{
				Label: "users'",
				UOM: "%",
				Value: 2,
				Warning: 3,
				Critical: 7,
				Min: 0,
				Max: math.NaN(),
			},
		},
		err: "could not parse perfdata: invalid format",
	},
	{
		in: "users2%;3;7;0",
		out: []*Perfdata{
			&Perfdata{
				Label: "users",
				UOM: "",
				Value: 2,
				Warning: 3,
				Critical: 7,
				Min: 0,
				Max: math.NaN(),
			},
		},
		err: "could not parse perfdata: no value found",
	},
	{
		in: " ",
		out: nil,
		err: "could not split perfdata:  ",
	},
	{
		in: "users=2;;;;;",
		out: nil,
		err: "could not parse perfdata: invalid perfdata string",
	},
	{
		in: "users=2.4.5",
		out: []*Perfdata{
			&Perfdata{
			},
		},
		err: "could not parse perfdata: strconv.ParseFloat: parsing \"2.4.5\": invalid syntax",
	},
	{
		in: "=2,34",
		out: []*Perfdata{
			&Perfdata{
			},
		},
		err: "could not parse perfdata: invalid label",
	},
	{
		in: "users=",
		out: []*Perfdata{
			&Perfdata{
			},
		},
		err: "could not parse perfdata: missing number",
	},
	{
		in: "users=1;2,4,4",
		out: []*Perfdata{
		},
		err: "could not parse perfdata: strconv.ParseFloat: parsing \"2.4.4\": invalid syntax",
	},
	{
		in: "users=1;;2,4,4",
		out: []*Perfdata{
		},
		err: "could not parse perfdata: strconv.ParseFloat: parsing \"2.4.4\": invalid syntax",
	},
	{
		in: "users=1;;;2,4,4",
		out: []*Perfdata{
		},
		err: "could not parse perfdata: strconv.ParseFloat: parsing \"2.4.4\": invalid syntax",
	},
	{
		in: "users=1;;;;2,4,4",
		out: []*Perfdata{
		},
		err: "could not parse perfdata: strconv.ParseFloat: parsing \"2.4.4\": invalid syntax",
	},
}

func testPerfdata(in string, out []*Perfdata) error {
	result, err := ParsePerfdata(in)
	if err != nil {
		return err
	}
	if len(result) != len(out) {
		return fmt.Errorf("number of parsed perfdata %d != expected %d", len(result), len(out))
	}
	for i := 0; i < len(result); i++ {
		if out[i].Label != result[i].Label {
			return fmt.Errorf("Label is not equal")
		}
		if out[i].Value != result[i].Value {
			return fmt.Errorf("Value is not equal")
		}
		if out[i].UOM != result[i].UOM {
			return fmt.Errorf("UOM is not equal")
		}
		if !(out[i].Warning == result[i].Warning || (math.IsNaN(out[i].Warning) && math.IsNaN(result[i].Warning))) {
			return fmt.Errorf("Warning is not equal")
		}
		if !(out[i].Critical == result[i].Critical || (math.IsNaN(out[i].Critical) && math.IsNaN(result[i].Critical))) {
			return fmt.Errorf("Critical is not equal")
		}
		if !(out[i].Min == result[i].Min || (math.IsNaN(out[i].Min) && math.IsNaN(result[i].Min))) {
			return fmt.Errorf("Min is not equal")
		}
		if !(out[i].Max == result[i].Max || (math.IsNaN(out[i].Max) && math.IsNaN(result[i].Max))) {
			return fmt.Errorf("Max is not equal")
		}
	}
	return nil
}

func TestParsePerfdata(t *testing.T) {
	for _, pftest := range testPerfdataSetOk {
		if err := testPerfdata(pftest.in, pftest.out); err != nil {
			t.Errorf("in: %s err: %s", pftest.in, err)
		}
	}
	for _, pftest := range testPerfdataSetFail {
		err := testPerfdata(pftest.in, pftest.out)
		if err == nil {
			t.Errorf("in: %s expected error \"%s\" but it passed", pftest.in, pftest.err)
		} else {
			if err.Error() != pftest.err {
				t.Errorf("in: %s expected error \"%s\" but got \"%s\"", pftest.in, pftest.err, err)
			}
		}
	}
}

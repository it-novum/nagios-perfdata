package perfdata

import (
	"math"
	"strconv"
	"strings"
	"fmt"
	"regexp"
)

// Perfdata contains all fields for a single nagios perfdata string
// the optional values Min, Max, Warning and Critical
// are set to math.NaN() if not specified in the string
type Perfdata struct {
	Label string
	Value float64
	UOM string
	Min float64
	Max float64
	Warning float64
	Critical float64
}

var perfdataSplitRegex = regexp.MustCompile(`(('[^']+')?[^\s]+)`)
var perfdataLabelRegex = regexp.MustCompile(`^(('[^']+')?([^='])*)`)
var perfdataValueRegex = regexp.MustCompile(`^([\d\.,]+)`)

func parseFloat(val string) (float64, error) {
	return strconv.ParseFloat(strings.Replace(val, ",", ".", -1), 64)
}

func perfdataParseValue(valueStr string) (*Perfdata, error) {
	var err error
	pd := &Perfdata{
		Label: "",
		Value: math.NaN(),
		UOM: "",
		Warning: math.NaN(),
		Critical: math.NaN(),
		Min: math.NaN(),
		Max: math.NaN(),
	}
	data := strings.Split(valueStr, ";")
	dataLen := len(data)
	if dataLen < 1 || dataLen > 5{
		return nil, fmt.Errorf("invalid perfdata string")
	}
	// 'label'=value[UOM];[warn];[crit];[min];[max]
	if dataLen == 5 && data[4] != "" {
		if pd.Max, err = parseFloat(data[4]); err != nil {
			return nil, err
		}
	}
	if dataLen >= 4 && data[3] != "" {
		if pd.Min, err = parseFloat(data[3]); err != nil {
			return nil, err
		}
	}
	if dataLen >= 3 && data[2] != "" {
		if pd.Critical, err = parseFloat(data[2]); err != nil {
			return nil, err
		}
	}
	if dataLen >= 2 && data[1] != "" {
		if pd.Warning, err = parseFloat(data[1]); err != nil {
			return nil, err
		}
	}
	pd.Label = perfdataLabelRegex.FindString(data[0])
	if pd.Label == "" {
		return nil, fmt.Errorf("invalid label")
	}
	if len(pd.Label) == len(data[0]) {
		return nil, fmt.Errorf("no value found")
	}
	if data[0][len(pd.Label)] != '=' {
		return nil, fmt.Errorf("invalid format")
	}
	valueWithUnit := data[0][len(pd.Label)+1:]
	rawval := perfdataValueRegex.FindString(valueWithUnit)
	if rawval == "" {
		return nil, fmt.Errorf("missing number")
	}
	if pd.Value, err = parseFloat(rawval); err != nil {
		return nil, err
	}

	if len(rawval) == len(valueWithUnit) {
		pd.UOM = ""
	} else {
		pd.UOM = valueWithUnit[len(rawval):]
	}

	return pd, nil
}

// ParsePerfdata splits a string into an array of *Perfdata values
func ParsePerfdata(perfdata string) ([]*Perfdata, error) {
	var err error
	valueStrings := perfdataSplitRegex.FindAllString(perfdata, -1)
	if len(valueStrings) == 0 {
		return nil, fmt.Errorf("could not split perfdata: %s", perfdata)
	}
	values := make([]*Perfdata, len(valueStrings))
	for i, valueStr := range valueStrings {
		if values[i], err = perfdataParseValue(valueStr); err != nil {
			return nil, fmt.Errorf("could not parse perfdata: %s", err)
		}
	}
	return values, nil
}

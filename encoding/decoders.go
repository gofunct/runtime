package encoding

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"strings"
)

// DefaultDecoders contains the default list of decoders per MIME type.
var DefaultDecoders = DecoderGroup{
	"xml":  DecoderMakerFunc(func(r io.Reader) Decoder { return xml.NewDecoder(r) }),
	"json": DecoderMakerFunc(func(r io.Reader) Decoder { return json.NewDecoder(r) }),
	"yaml": DecoderMakerFunc(func(r io.Reader) Decoder { return &yamlDecoder{r} }),
}

type (
	// A Decoder decodes data into v.
	Decoder interface {
		Decode(v interface{}) error
	}

	// A DecoderGroup maps MIME types to DecoderMakers.
	DecoderGroup map[string]DecoderMaker

	// A DecoderMaker creates and returns a new Decoder.
	DecoderMaker interface {
		NewDecoder(r io.Reader) Decoder
	}

	// DecoderMakerFunc is an adapter for creating DecoderMakers
	// from functions.
	DecoderMakerFunc func(r io.Reader) Decoder
)

// NewDecoder implements the DecoderMaker interface.
func (f DecoderMakerFunc) NewDecoder(r io.Reader) Decoder {
	return f(r)
}

type yamlDecoder struct {
	r io.Reader
}

func (yd *yamlDecoder) Decode(v interface{}) error {
	b, err := ioutil.ReadAll(yd.r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

func ReadAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func ReadAsMap(val string) (map[string]string, error) {
	var newMap = make(map[string]string)
	Slice, err := ReadAsCSV(val)
	if err != nil {
		return nil, err
	}
	for _, str := range Slice { // iterating over each tab in the csv
		//map k:v are seperated by either = or : and then a comma
		strings.TrimSpace(str)
		if strings.Contains(str, "=") {
			newSlice := strings.Split(str, "=")
			newMap[newSlice[0]] = newSlice[1]
		}
		if strings.Contains(str, ":") {
			newSlice := strings.Split(str, ":")
			newMap[newSlice[0]] = newSlice[1]
		}
	}
	if newMap == nil {
		return nil, errors.New("cannot conver string to map[string]string- detected a nil map output")
	}
	return newMap, nil
}

// toJson encodes an item into a JSON string
func ToJson(v interface{}) string {
	output, _ := json.Marshal(v)
	return string(output)
}

// toPrettyJson encodes an item into a pretty (indented) JSON string
func ToPrettyJson(v interface{}) string {
	output, _ := json.MarshalIndent(v, "", "  ")
	return string(output)
}

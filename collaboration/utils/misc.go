package utils

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func UnmarshalStrict(in []byte, out interface{}) (err error) {
	knownFieldsDecoder := yaml.NewDecoder(bytes.NewReader(in))
	knownFieldsDecoder.KnownFields(true)
	return knownFieldsDecoder.Decode(out)
}
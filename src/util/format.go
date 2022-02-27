package util

type ResultData interface {
	ToJson() string
	ToString() string
	ToFormat() string
}

func Format(d ResultData, isFormat, isJson bool) string {
	if isFormat {
		return d.ToFormat()
	}
	if isJson {
		return d.ToJson()
	}
	return d.ToString()
}

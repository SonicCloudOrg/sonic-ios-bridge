package tool

type Data interface {
	ToJson() string
	ToString() string
	ToFormat() string
}

func Format(d Data, isFormat, isJson bool) string {
	if isFormat {
		return d.ToFormat()
	}
	if isJson {
		return d.ToJson()
	}
	return d.ToString()
}

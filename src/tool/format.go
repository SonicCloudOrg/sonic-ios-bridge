package tool

type Data interface {
	ToJson() string
	ToString() string
}

func Format(d Data, s string) string {
	switch s {
	case "json":
		return d.ToJson()
	case "string":
		return d.ToString()
	default:
		return ""
	}
}

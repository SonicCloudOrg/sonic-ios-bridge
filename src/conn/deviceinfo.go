package conn

type GetValueRequest struct {
	Label   string
	Key     string `plist:"Key,omitempty"`
	Request string `plist:"Request"`
	Domain  string `plist:"Domain,omitempty"`
	Value   string `plist:"Value,omitempty"`
}

func NewGetValue(domain string, key string) GetValueRequest {
	data := GetValueRequest{
		Label:   BundleId,
		Domain:  domain,
		Key:     key,
		Request: "GetValue",
	}
	return data
}

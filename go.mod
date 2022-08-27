module github.com/SonicCloudOrg/sonic-ios-bridge

go 1.18

require (
	github.com/electricbubble/gidevice v0.6.2
	github.com/google/uuid v1.3.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.4.0
	github.com/valyala/fastjson v1.6.3
	howett.net/plist v1.0.0
)

replace github.com/electricbubble/gidevice v0.6.2 => github.com/SonicCloudOrg/sonic-gidevice v0.0.0-20220827051900-b54b5c523e71

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

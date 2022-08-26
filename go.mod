module github.com/SonicCloudOrg/sonic-ios-bridge

go 1.18

require (
	github.com/electricbubble/gidevice v0.6.2
	github.com/mitchellh/mapstructure v1.5.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.4.0
)
// todo 修改为序列化信息优化的版本
replace github.com/electricbubble/gidevice v0.6.2 => github.com/SonicCloudOrg/sonic-gidevice v0.0.0-20220809142714-8bef4cc76426

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	howett.net/plist v1.0.0 // indirect
)

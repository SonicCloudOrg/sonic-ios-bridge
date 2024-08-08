<p align="center">
  <img width="80px" src="https://raw.githubusercontent.com/SonicCloudOrg/sonic-server/main/logo.png">
</p>
<p align="center">🎉Bridge of iOS Devices by usbmuxd</p>
<p align="center">
  <span>English |</span>
  <a href="https://github.com/SonicCloudOrg/sonic-ios-bridge/blob/main/README_CN.md">  
     简体中文
  </a>
</p>
<p align="center">
  <a href="#">  
    <img src="https://img.shields.io/github/v/release/SonicCloudOrg/sonic-ios-bridge?include_prereleases">
  </a>
  <a href="#">  
    <img src="https://img.shields.io/github/downloads/SonicCloudOrg/sonic-ios-bridge/total">
  </a>
<a href="https://app.fossa.com/projects/git%2Bgithub.com%2FSonicCloudOrg%2Fsonic-ios-bridge?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2FSonicCloudOrg%2Fsonic-ios-bridge.svg?type=shield"/></a>
  <a href="#">  
    <img src="https://img.shields.io/github/go-mod/go-version/SonicCloudOrg/sonic-ios-bridge">
  </a>
</p>

## Document
[Sonic Official Website](https://soniccloudorg.github.io/sib/re-sib.html)

## Usage

#### 1. Download
[click here](https://github.com/SonicCloudOrg/sonic-ios-bridge/releases)
#### 2. execute shell (No need for windows)
```
sudo chmod 777 ./sib && ./sib version
```
#### 3. Add sib to your PATH
Finish!


## Function
You should mount before use it.
```
sib mount
```
then
```
sib run wda -b your.wda.bundleId
sib run xctest -b your.wda.bundleId
sib remote share
sib remote connect --host 192.168.1.1
sib app list
sib app launch
sib devices listen
sib app uninstall
sib screenshoot
sib ps
sib crash
sib location
sib oritation
sib battery
sib info
...
```
👉[ (Recommend) Click Here to Get More!](https://soniccloudorg.github.io/sib/re-sib.html)


## Sponsors

Thank you to all our sponsors!

[<img src="https://ceshiren.com/uploads/default/original/3X/7/0/70299922296e93e2dcab223153a928c4bfb27df9.jpeg" alt="霍格沃兹测试开发学社" width="500">](https://qrcode.testing-studio.com/f?from=sonic&url=https://ceshiren.com)

> [霍格沃兹测试开发学社](https://qrcode.testing-studio.com/f?from=sonic&url=https://ceshiren.com)是业界领先的测试开发技术高端教育品牌，隶属于[测吧（北京）科技有限公司](http://qrcode.testing-studio.com/f?from=sonic&url=https://www.testing-studio.com) 。学院课程由一线大厂测试经理与资深测试开发专家参与研发，实战驱动。课程涵盖 web/app 自动化测试、接口测试、性能测试、安全测试、持续集成/持续交付/DevOps，测试左移&右移、精准测试、测试平台开发、测试管理等内容，帮助测试工程师实现测试开发技术转型。通过优秀的学社制度（奖学金、内推返学费、行业竞赛等多种方式）来实现学员、学社及用人企业的三方共赢。[进入测试开发技术能力测评!](https://qrcode.testing-studio.com/f?from=sonic&url=https://ceshiren.com/t/topic/14940)

## Thanks

- [https://github.com/electricbubble/gidevice](https://github.com/electricbubble/gidevice)
- [https://github.com/libimobiledevice/libimobiledevice](https://github.com/libimobiledevice/libimobiledevice)
- [https://github.com/danielpaulus/go-ios](https://github.com/danielpaulus/go-ios)

## LICENSE

[License](LICENSE)


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FSonicCloudOrg%2Fsonic-ios-bridge.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FSonicCloudOrg%2Fsonic-ios-bridge?ref=badge_large)

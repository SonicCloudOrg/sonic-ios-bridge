<p align="center">
  <img width="80px" src="https://raw.githubusercontent.com/SonicCloudOrg/sonic-server/main/logo.png">
</p>
<p align="center">🎉基于usbmuxd的iOS调试工具</p>
<p align="center">
  <a href="https://github.com/SonicCloudOrg/sonic-ios-bridge/blob/main/README.md">  
    English
  </a>
  <span>| 简体中文</span>
</p>
<p align="center">
  <a href="#">  
    <img src="https://img.shields.io/github/v/release/SonicCloudOrg/sonic-ios-bridge?include_prereleases">
  </a>
   <a href="#">  
    <img src="https://img.shields.io/github/downloads/SonicCloudOrg/sonic-ios-bridge/total">
  </a>
  <a href="#">  
    <img src="https://img.shields.io/github/go-mod/go-version/SonicCloudOrg/sonic-ios-bridge">
  </a>
</p>

### 官方文档
[Sonic Official Website](https://sonic-cloud.cn/sib/re-sib.html)

## 使用方法

#### 1. 下载
[点击这里](https://github.com/SonicCloudOrg/sonic-ios-bridge/releases)

#### 2. 执行指令 (windows不需要)
```
sudo chmod 777 ./sib && ./sib version
```

#### 3. 添加sib路径到本机PATH
完成！

## 功能
使用前应该要先mount
```
sib mount
```
然后
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
👉[ (推荐) 点击这里查看更多!](https://sonic-cloud.cn/sib/re-sib.html)

## 赞助商

感谢所有赞助商！

[<img src="https://ceshiren.com/uploads/default/original/3X/7/0/70299922296e93e2dcab223153a928c4bfb27df9.jpeg" alt="霍格沃兹测试开发学社" width="500">](https://qrcode.testing-studio.com/f?from=sonic&url=https://ceshiren.com)

> [霍格沃兹测试开发学社](https://qrcode.testing-studio.com/f?from=sonic&url=https://ceshiren.com)是业界领先的测试开发技术高端教育品牌，隶属于[测吧（北京）科技有限公司](http://qrcode.testing-studio.com/f?from=sonic&url=https://www.testing-studio.com) 。学院课程由一线大厂测试经理与资深测试开发专家参与研发，实战驱动。课程涵盖 web/app 自动化测试、接口测试、性能测试、安全测试、持续集成/持续交付/DevOps，测试左移&右移、精准测试、测试平台开发、测试管理等内容，帮助测试工程师实现测试开发技术转型。通过优秀的学社制度（奖学金、内推返学费、行业竞赛等多种方式）来实现学员、学社及用人企业的三方共赢。[进入测试开发技术能力测评!](https://qrcode.testing-studio.com/f?from=sonic&url=https://ceshiren.com/t/topic/14940)
 
## 感谢

- [https://github.com/electricbubble/gidevice](https://github.com/electricbubble/gidevice)
- [https://github.com/libimobiledevice/libimobiledevice](https://github.com/libimobiledevice/libimobiledevice)
- [https://github.com/danielpaulus/go-ios](https://github.com/danielpaulus/go-ios)

## 开源许可协议

[License](LICENSE)

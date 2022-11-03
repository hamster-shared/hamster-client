## wails 开发

项目地址： https://github.com/hamster-shared/hamster-client.git

开发基础分支： develop  使用[github flow](https://docs.github.com/cn/get-started/quickstart/github-flow)

wsUrl: `ws://183.66.65.207:49944`

页面原型： http://183.66.65.205:32453/#id=1xx5dz&p=%E9%83%A8%E7%BD%B2%E5%BA%94%E7%94%A8&g=1
效果图： https://lanhuapp.com/url/j2WoI-AjWeaF


### 1. 项目编译

go version: 1.17+

```bash

#1. 设置代理  windows 参见：https://goproxy.cn/
export GO111MODULE=on
export GOPROXY=https://goproxy.cn


#2. 下载wails 命令行
go install github.com/wailsapp/wails/v2/cmd/wails@latest

#3. 检查wails 环境完整
wails doctor

#4. 以dev 方式启动环境
wails dev


#5. 网页访问
http://localhost:8085


```


### 2. 前后端调用



前端： window.go.app.<模块名>.<方法名>（...params).then(res => {
do_some_response(...)
})


参考： frontend/src/views/setting/index.vue
 
      frontend/src/wailsjs/go/app/Setting.d.ts

examplse:
```javascript
window.go.app.Setting.Setting(state.settingForm.publicKey,state.settingForm.wsUrl).then(res => { ... } )

window.go.app.Wallet.DeleteWallet();
```


后端： /app/*.go

example :
```go
func (s *Setting) GetSetting() (*Config, error) { ... }

func (s *Setting) Setting(publicKey string, wsUrl string) (bool, error) { ... }

```

### 3. 特别注意

1. 开源纯英文项目， 可以有国际化，但是git commit log 不能有中文

2. 后端接口，找孙建国
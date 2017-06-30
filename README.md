# goImg

## 简介
goImg是一个使用golang语言编写的图片服务器，支持多种图片格式和多种上传方式以及图片尺寸缩放。
## 特点
* 图片格式支持jpeg、gif、png、webp格式。
* 文件存储目录采用md5算法生成。
* 图片缩放只支持jpeg格式。
## 安装
go get  github.com/nymbian/goImg
## 配置文件
conf.json
```JSON
{
	"ListenAddr":":10086",      <-监听地址
	"Storage":"/data/image/"      <-存储位置
}
```

## 使用方法
### 上传图片方法1
#### 文件上传方式
POST /upload
表单参数:
file file类型,要上传的图片
返回值:
{图片ID}

### 上传图片方法2
#### 图片url方式
POST /url
表单参数:
url 要上传的图片的远程地址
返回值:
{图片ID}

### 上传图片方法3
#### base64方式 
POST /base64
表单参数:
base64 要上传的图片的base64值
返回值:
{图片ID}

### 获取图片
#### 原图
GET /{图片ID}
返回值:
{图片文件}
#### resize
GET /{图片ID}_{宽}x{高}
返回值:
{图片文件}

### 第三方组件
* github.com/gorilla/mux
* github.com/nfnt/resize

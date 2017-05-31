# goImg

## 简介
goImg是一个简单的图像服务器。
## 特点
* 图片格式支持jpeg、gif、png、webp格式。
* 文件存储目录采用md5算法生成。
* 图片缩放只支持jpeg格式。
##安装
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
### 上传图片
POST /
表单参数:
uploadfile file类型,要上传的图片
返回值:
{图片ID}

### 获取图片
原图
GET /{图片ID}
返回值:
{图片文件}
resize原图
GET /{图片ID}_{宽}x{高}
返回值:
{图片文件}

### 第三方组件
* github.com/gorilla/mux
* github.com/nfnt/resize

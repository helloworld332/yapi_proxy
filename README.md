# yapi_proxy
一个实用的测试代理服务，与fiddler类似，fiddler的核心是专注抓包进行单请求的跟踪分析，而本代理的核心专注于接口代理的高度配置化和自动化，
以满足测试人员灵活的进行接口响应方的配置，yapi和rap都是目前开源的接口数据Mock工具，搭配上它们系统上编辑的接口Mock数据，
通过插拔和替换接口配置文件，可以灵活的、高度自动化整个测试流程。

yapi和rap的官方网址如下，建议使用yapi，功能更齐全，团队支持力度更大。
- yapi: http://yapi.demo.qunar.com/
- rap: http://rapapi.org/org/index.do


# 安装
- 安装好golang开发环境，开发环境安装配置方法: https://www.jianshu.com/p/eb35a47a157e
- 通过go get -v github.com/helloworld332/yapi_proxy/ 获取最新代码
- 进入代码的server目录，执行go build 生成可执行文件
- 具体实用方法请看代码的doc目录下的使用手册

# 项目依赖
- 依赖goproxy[https://github.com/elazarl/goproxy] ， 它是我们做代理的底层

# 核心代码目录
- server: 代理的核心代码
- conf: 配置文件
- yapi_tools: 在自己本机安装yapi的情况下，可以连接到yapi的mongodb数据库进行项目的复制，删除等操作，这是yapi官方不提供的功能。

# 项目名称

基于openwechat框架实现的feybot机器人
实现功能
1. 微信群消息+个人消息的自动回复
2. 接入GPT-3.5以及通义千问AI模型
3. 实现了插件注册功能(便于扩展性)
4. 实现了热调试控制台(日志热调试+插件热调试)

## 注意
写代码小白，很多地方设计和代码都没有规范
## 开始

### 下载代码到本地
通过git clone
```bash
$ git@github.com:CoderStarfeyfey/feybot.git
```
或者
Download ZIP下载源码

## 项目架构
1.单线程+协程

该方案比较简单,整个程序只采用简单的单线程，不会涉及到多个进程以及多个线程之间交互的场景，直接在openwechat框架提供的bot对象的结构体中注册消息处理的回调函数，openwechat框架已经实现了消息的长轮询检查(代码是同步这里会被阻塞)，每当收到消息后就对收到的消息都做实现定义好的回调函数。除此之外，定时任务执行也是比较关键的，像每日定时推送天气，新闻或者定时清理后台文件以及定时清理数据库都是比较重要的，在webot的主进程启动的时候开启协程，让协程执行这个定时任务即可。因此这个feybot机器人开发是比较简单的单线程+协程而已。

此外整个项目为了实现可扩展的灵活性，项目是基于插件化实现的，具体的处理消息的有不同的插件，比如说有根据配置文件定义好的回复来回复的插件以及调用chatgpt的接口来回复的插件等等，和C语言的的dll一样，运行程序的时候会执行指定目录下的go的init函数作为插件的入口函数，在入口函数中要将要实现业务逻辑的函数添加到一个map中，这个map是全局可见的map专门用来为插件提供注册入口函数的功能。后续在机器人的处理逻辑中会便利这个map去执行不同插件的消息处理的逻辑。

项目还实现了控制台的热调试功能，控制台输入注册好的命令即可热调试。这里是基于IPC实现一个控制台进程，控制台进程利用套接字或者管道或其余的IPC通信的方法去和和主进程进行通信，当然这个主进程一定是开启一个协程，协程和这个控制台通信并且去接受控制条输入的命令，并执行对应的函数输出字符串到控制台上显示。这里也是要用到注册调试命令的思想。大概思路和插件注册差不多，这里做的就是在插件注册成功后，在一个全局的map中注册一个调试命令以及对应的处理函数的回调，这里要求回调的签名是一样的。


## 使用

用户只需要根据config.go里定义的结构体来配置配置文件即可。本项目有两个配置文件分别归档在config/feybot.config以及normalReplyConfig.config，其中normalReplyConfig.config是回复指定消息的配置，feybot.config是机器人最重要的配置，下面给出一个示例，当然token被隐藏掉了，用户需要自己去平台申请一个token。
```text
{
  "botName" : "feybot",
  "dataDir" : "./botlog",
  "tyqwToken" : "token",
  "chatGPTToken" : "token"
  "defalutReplayMsg" : {
    "coolManMsg" : [
      "Cool"
    ],
    "errorMsg" : "消息错误！",
    "thankMsg" : {
      "thankMsgText" : "帮助你是我应该的事情，不用说感谢的话，拿出你的行动",
      "ThankMsgPic" : "receipt/qr.jpg"
    }
  },
  "groupWhiteList" : [
    "测试群1",
    "测试群2"
  ],
  "feyLogConfig" : {
    "maxsize": 5,
    "maxbackups": 3,
    "maxage": 30,
    "compress": true
  },
  "feature" : {
    "gpt自动回复" : {
      "enable" : true,
      "entryFunctionName" : "GPTAutoReply",
      "enableGroupWbList" : {
        "测试群1" : true,
        "测试群2" : true
      }
    },
    "摸鱼日历" : {
      "enable" : true,
      "entryFunctionName" : "HolidaysQuery",
      "enableGroupWbList" : {
        "测试群1" : true,
        "测试群2" : true
      }
    },
    "每日一题" : {
      "enable" : true,
      "entryFunctionName" : "LeetcodeReply",
      "enableGroupWbList" : {
        "测试群1" : true,
        "测试群2" : true
      }
    },
    "高级情话" : {
          "enable" : true,
          "entryFunctionName" : "LoveWords",
          "enableGroupWbList" : {
            "测试群1" : true,
            "测试群2" : true
          }
        }
  }
}
```
# LAN-chat
> 局域网聊天工具

## 需求说明：

​    实现应用于局域网的内部通信工具。

## 功能点：

1. 安装该工具的局域网主机间可进行文字通信；`websocket`
2. 安装该工具的局域网主机间可发送语音、不超过50M的视频；`websocket`
3. 安装该工具的局域网主机间可发送不超过5M的非视频文件；`websocket`
4. 安装该工具的局域网主机间可按照用户配置进行文件类型过滤；
5. 该工具可保留用户通信记录；
6. 用户可删除本地的通信记录；
7. ~~工具支持不同的用户以用户名、密码登录；~~（每个设备会保留自己的指纹作为唯一表示符）
8. ~~服务端保留用户的认证信息（用户名、密码）；~~（没有必要）
9. 服务端保留用户的通信记录；

# 接口文档

## 登录模块

### 1. 登录接口

- 请求方式：POST

- URL：`http://localhost:80/login`

- 请求体

  - ```json
    {
        "userName": "江键翔",
        "passWord"： "12345678"
    }
    ```

    - 参数说明：
      - "userName"：用户名
      -  "passWord"：密码

- 返回结果：

  - 成功：返回成功消息

    - ```json
      {
          "code": 1,
          "data": {
             "message":  "用户登录成功",
             "token": "sdsadsadadddsfds"
      }
      ```

  - 失败：返回错误信息

    - ```json
      {
          "code": 0,
          "data": {
             "message":  "用户登录失败"
      }
      ```

      * 参数说明：
        * "code"：1表示成功，0表示失败
        * "data"：信息
          * "message"：说明
          * "token"：包含`uId`等信息

### 2. 注册接口

- 请求方式：POST

- URL：`http://localhost:80/user/register`

- 请求体

  - ```json
    {
        "userName": "江键翔"
        "passWord"： "12345678"
        "email": "1364054111@qq.com"
        "verifyNum": "12345"
    }
    ```

    - 参数说明：
      - "userName"：用户名
      -  "passWord"：密码
      - "email"：邮箱
      - "verifyNum"：验证码

- 返回结果：

  - 成功：返回成功消息

    - ```json
      {
          "code": 1,
          "data": {
             "message":  "用户注册成功"
      }
      ```

  - 失败：返回错误信息

    - ```json
      {
          "code": 0,
          "data": {
             "message":  "用户注册失败"
      }
      ```

      * 参数说明：
        * "code"：1表示成功，0表示失败
        * "data"：信息
          * "message"：说明

### 3. 修改密码模块

#### 3.1 **生成验证码**

- 请求方式：POST

- URL：`http://localhost:80/user/chgPsd/generate`

- 请求体

  - ```json
    {
        "email": "1364054111@qq.com"
    }
    ```

    - 参数说明：
      - "email"：邮箱
      - "verifyNum"：验证码

- 返回结果：

  - 成功：返回成功消息

    - ```json
      {
          "code": 1,
          "data": {
             "message":  "邮箱存在，请查看验证码",
             "verifyNum"："56327"
      }
      ```

  - 失败：返回错误信息

    - ```json
      {
          "code": 0,
          "data": {
             "message":  "邮箱不存在，请确保输入的邮箱正确!"
      }
      ```

      * 参数说明：
        * "code"：1表示成功，0表示失败
        * "data"：信息
          * "message"：说明
          * "verifyNum"：验证码

#### 3.2 **提交新密码**

- 请求方式：POST

- URL：`http://localhost:80/user/chgPsd/update`

- 请求体

  - ```
     {
         "userName": "江键翔"
         "passWord"： "12345678"
     }
    ```

    - 参数说明：
      - "userName"：用户名
      -  "passWord"：密码

- 返回结果：

  - 成功：返回成功消息

    - ```
       {
           "code": 1,
           "data": {
              "message":  "修改密码成功"
       }
      ```

  - 失败：返回错误信息

    - ```
       {
           "code": 0,
           "data": {
              "message":  "修改密码失败"
       }
      ```

      - 参数说明：
        - "code"：1表示成功，0表示失败
        - "data"：信息
          - "message"：说明

## 通信模块

### 消息发送模块

#### 1. 发送文本消息接口

- 请求方式：POST

- URL：`http://localhost:80/user/sendMesg`

- 请求头

  - ```json
     {
         "token": "osgfjdhgfdsofdsdfkf3w3"
     }
    ```

- 请求体

  - ```json
     {
         "message": "hello world！",
         "recipient":
     }
    ```

  - 参数说明：

    - 接收者用户名（recipient）
    - 消息内容（message）
    - 身份验证令牌（token）

- 返回结果：

  - 成功：返回成功消息

    - ```json
       {
           "code": 1,
           "data": {
              "message":  "修改密码成功"
       }
      ```

  - 失败：返回错误信息

    - ```json
       {
           "code": 0,
           "data": {
              "message":  "修改密码失败"
       }
      ```

      - 参数说明：
        - "code"：1表示成功，0表示失败
        - "data"：信息
          - "message"：说明

#### 2. 发送语音消息接口

- 请求方式：POST

- URL：`http://localhost:80/user/sendVoice`

- 参数：

  - 接收者用户名（recipient）
  - 消息内容（message）
  - 身份验证令牌（token）

- 返回结果：

  - 成功：返回成功消息

    - ```JSON
      {
          "code": 1,
          "data": {
             "message":  "查询个人信息成功，已返回全部个人信息",
             "userInfo": userInfo,
      }
      ```

  - 失败：返回错误信息

#### 3. 发送视频消息接口

- 请求方式：POST
- URL：`http://localhost:8000/sendVideo`
- 参数：
  - 接收者用户名（recipient）
  - 语音文件（voice）
  - 身份验证令牌（token）
- 返回结果：
  - 成功：返回成功消息
  - 失败：返回错误信息

#### 4. 发送文件接口

- 请求方式：POST
- URL：`http://localhost:8000/sendFile`
- 参数：
  - 接收者用户名（recipient）
  - 视频文件（video）
  - 身份验证令牌（token）
- 返回结果：
  - 成功：返回成功消息
  - 失败：返回错误信息

### 消息记录模块

#### 1. 获取通信记录接口

- 请求方式：POST
- URL：`http://localhost:80/user/getRecord`
- 参数：
  - 接收者用户名（recipient）
  - 文件（file）
  - 身份验证令牌（token）
- 返回结果：
  - 成功：返回成功消息
  - 失败：返回错误信息

#### 2. 删除通信记录接口

- 请求方式：DELETE
- URL：`http://localhost:80/user/deleteRecord`
- 参数：
  - 消息ID（message_id）
  - 身份验证令牌（token）
- 返回结果：
  - 成功：返回成功消息
  - 失败：返回错误信息



```json
{
  "type": "composite",
  "content": [
    {
      "partType": "text",
      "data": "Here is a summary of the content."
    },
    {
      "partType": "image",
      "data": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAU...（数据省略）"
    },
    {
      "partType": "video",
      "data": "data:video/mp4;base64,AAAAIGZ0eXBtcDQyAAAAAG1wN...（数据省略）"
    }
  ],
  "timestamp": 1618158147,
  "sender": "user123"
}

```



# 请求（Request）

## 请求`JSON`结构

```json
{
    "data": {
        "type": "request",
        "timestamp": "时间戳",
        "from": "发送方标识",
        "to": "接收方标识"
    },
    "content": [
        {
            "partType": "内容部分类型",
            "data": "具体数据"
        }
    ]
}
```

## 字段说明

### data字段

- data: 包含了消息的基本信息
  - type: 操作类型，此处为send表示发送操作。
  - `timestamp`: 时间戳，同时作为唯一标识符，遵循ISO 8601格式。
  - from: 发送方的唯一标识。
  - to: 接收方的唯一标识。

### content字段

- content: 消息内容数组，可以包含多种类型的内容（如文本、图片、视频和音频）。
  - `partType`: 内容部分的类型，可以是text、image、video或audio。
  - data: 具体的数据内容。对于文本，直接是字符串；对于图片、视频和音频，则是相应的base64编码数据。

## 示例

```json
{
    "data": {
        "type": "request",
        "timestamp": "2024-02-18T12:34:56Z"， // 同时作为唯一标志符
        "from": "user123",
        "to": "user456",
    },
    "content": [
        {
            "partType": "text",
            "data": "Here is a summary of the content."
        },
        {
            "partType": "image",
            "data": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAU"
        },
        {
            "partType": "video",
            "data": "data:video/mp4;base64,AAAAIGZ0eXBtcDQyAAAAAG1wN"
        },
        {
            "partType": "audio",
            "data": "data:audio/mp3;base64,SUQzBAAAAAAAI1RTU0UAAA"
        }
    ]
}
```

# 返回（Response）

## 请求`JSON`结构

```json
{
    "data": {
        "type": "response",
        "timestamp": "时间戳",
        "from": "响应方标识",
        "to": "请求方标识"
    },
    "content": {
        "status": "操作状态",
        "message": "状态消息"
    }
}
```

### 字段说明

#### data字段

- `type`: 操作类型，此处为`response`表示这是一个响应操作。
- `timestamp`: 时间戳，同时作为该响应的唯一标识符，遵循ISO 8601格式。这有助于确保数据的唯一性和顺序性，同时也方便日志记录和问题追踪。
- `from`: 响应方的唯一标识，表示这个响应是由谁发起的。
- `to`: 请求方的唯一标识，表示这个响应是发送给谁的。

#### content字段

- `status`: 操作的状态，例如`1`表示操作成功。`0`表示失败。
- `message`: 提供关于操作状态的额外信息或错误消息。例如，当`status`为`1`时，`message`可以是`"Message sent successfully."`，提供了更人性化的反馈。

### 示例

```json
{
    "data": {
        "type": "response",
        "timestamp": "2024-02-18T12:34:56Z"， // 同时作为唯一标志符
        "from": "user456",
        "to": "user123",
    },
    "content": {
        "status": 1,
        "message": "Message sent successfully."
    }
}
```







# 大问题

将 Qt 的代码直接嵌入到 Go 项目中并不是一个常见的做法，因为 Qt 主要是使用 C++ 开发的，而 Go 是一个不同的编程语言，它们之间的直接集成并不直接支持。不过，有几种方法可以实现 Qt 前端和 Go 后端之间的交互，以便用 Go 处理 Qt 中用户输入的文本：

### 1. 网络通信（推荐）

最常见的方式是通过网络通信，如之前提到的使用 WebSocket。Qt 客户端可以将用户输入的文本发送到 Go 服务器，服务器处理后可以返回结果。这种方法适用于客户端和服务器可能在不同机器上的情况，也适用于局域网或互联网应用。

### 2. 使用 Cgo 调用 C++ 代码

Go 提供了一个名为 `cgo` 的特性，允许 Go 程序调用 C 代码。由于 Qt 是用 C++ 编写的，你可以创建一个 C 风格的接口作为 Go 和 Qt 之间的桥梁。这需要编写一些 C（或C++）代码来封装 Qt 的功能，并通过 `extern "C"` 让它们可以被 Go 代码调用。这种方法比较复杂，涉及到跨语言编程和内存管理的问题，一般不推荐用于大型项目。

### 3. 使用本地进程间通信（IPC）

如果 Qt 客户端和 Go 服务器都在同一台机器上运行，你可以使用本地进程间通信（IPC）机制，如命名管道、共享内存或 UNIX 套接字等。Qt 客户端可以将数据写入一个共享的管道或套接字，Go 服务器则从中读取数据进行处理。这种方法适用于桌面应用程序或其他本地运行的软件。

### 4. 通过文件系统交互

另一种简单的方法是通过文件系统交互。Qt 客户端将用户输入写入一个特定的文件，Go 服务器监视该文件的变化，一旦发现变化就读取内容进行处理。这种方法实现简单，但效率不高，适用于对实时性要求不高的应用。

### 选择哪种方法？

- 对于大多数应用程序，特别是需要实时交云互或分布式系统的情况，**网络通信**是最灵活和最常用的解决方案。
- 如果你的应用程序运行在同一台机器上，并且对实时性要求不高，可以考虑使用**IPC**或**文件系统交互**。
- **Cgo**可以提供更紧密的集成，但需要更多的工作，并涉及复杂的内存管理和错误处理。

总的来说，选择哪种方法取决于你的具体需求、应用的架构以及你对性能和实时性的要求。在大多数情况下，使用网络通信（如 WebSocket）来在 Qt 前端和 Go 后端之间传递数据是最简单和最有效的方式。





**我还是把握不准。现在是这么个情况。在p2p架构下。**
**方案1.在程序内容写一个go的websocket服务端。它既要和qt的websocket客户端通信，也要和外部阶段通信。如果别人主动连接，go为服务端。如果我们要主动连接别人，go的websocket也要具备客户端的能力。**
**方案2.舍弃go作为中间桥梁。直接用qt搭建websocket客户端和服务端，直接与外部通信。请问我改选哪一种方案？**

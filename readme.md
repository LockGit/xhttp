## xhttp

### what is xhttp ?

xhttp ---> 将一个或者多个http api进行并行或串行编排后并将结果返回

### 配置定义

storage/demo下的文件夹表示项目名，每个项目下的routes.json为要编排的路由

| 字段                                  | 类型     | 说明                                               |
|-------------------------------------|--------|--------------------------------------------------|
| url                                 | string | 匹配url,指定值或/path/:id占位符或/path/some*通配符形式均可        |
| method                              | string | 当前项目url所支持的http请求                                |
| execType                            | string | parallel=并行，serial=串行                            |
| children                            | array  | 待编排的子api                                         |
| children[x].name                    | string | 子api的别名                                          |
| children[x].url                     | string | 子api的url                                         |
| children[x].method                  | string | 子api的请求方式                                        |
| children[x].params                  | array  | 子api的请求参数                                        |
| children[x].params[x].name          | string | 子api的请求参数名                                       |
| children[x].params[x].source        | string | 子api的请求参数取值来源                                    |
| children[x].params[x].default_value | string | 子api的请求参数默认值                                     |
| children[x].params[x].required      | string | 子api的请求参数是否必填                                    |

其中取值来源children[x].params[x].source 支持有以下：

| 字段       | 类型     | 说明                |
|----------|--------|-------------------|
| $.query  | string | 从http query中获取参数  |
| $.body   | string | 从http query中获取参数  |
| $.header | string | 从http header中获取参数 |
| $.cookie | string | 从http cookie中获取参数 |

### demo

```bash
go run main.go

curl -XPOST  http://127.0.0.1:8888/test?prod=env&from=pc&json=1 -H 'X-Project:hello' 

X-Project的值即为项目名,上面curl请求表示当前请求发到那个hello项目
```

服务接收到请求后，会匹配/test路由，/test下有2个子路由：

* 表示对children中的两个url（name=aaa,name=bbb）进行parallel(并行)聚合操作。
* 以children解释说明： 两个children中的url都配置了: https://www.baidu.com/sugrec
* 第2个children的get请求的参数有如下：
    * prod(从query中取值)
    * from(从header中取值)
    * name(从body中取值)
    * age(从cookie中取值)
    * other(从$.aaa【前一个name=aaa的api结果中取值】,仅在串行聚合中生效)

```json
{
  "url": "/test",
  "method": "post",
  "execType": "parallel",
  "children": [
    {
      "name": "aaa",
      "url": "https://www.baidu.com/sugrec",
      "method": "get",
      "params": [
        {
          "name": "prod",
          "source": "$.query",
          "default_value": "pc_his",
          "required": true
        },
        {
          "name": "from",
          "source": "$.header",
          "default_value": "pc",
          "required": true
        },
        {
          "name": "name",
          "source": "$.body",
          "default_value": "1",
          "required": true
        },
        {
          "name": "age",
          "source": "$.cookie",
          "default_value": "1",
          "required": true
        }
      ],
      "timeout": 3
    },
    {
      "name": "bbb",
      "url": "https://www.baidu.com/sugrec",
      "method": "get",
      "params": [
        {
          "name": "prod",
          "source": "$.query",
          "default_value": "pc_his",
          "required": true
        },
        {
          "name": "from",
          "source": "$.query",
          "default_value": "pc",
          "required": true
        },
        {
          "name": "name",
          "source": "$.query",
          "default_value": "1",
          "required": true
        },
        {
          "name": "age",
          "source": "$.cookie",
          "default_value": "1",
          "required": true
        },
        {
          "name": "other",
          "source": "$.aaa",
          "default_value": "test",
          "required": true
        }
      ],
      "timeout": 3
    }
  ]
}
```

```bash
curl -XPOST  http://127.0.0.1:8888/demo -H 'X-Project:hello'
```

```json
{
  "aaa": "{\"err_no\":0,\"errmsg\":\"\",\"queryid\":\"0x4a72d4bbf5c4bd\"}",
  "bbb": "{\"err_no\":0,\"errmsg\":\"\",\"queryid\":\"0x1acd194c0f99b0\"}"
}
```

### 注意

* children[x].url 的配置应该是确定的，否则可能会触发ssrf
* 用go写的速成的项目，用于测试目的， xhttp 还未经过大量的验证和实践
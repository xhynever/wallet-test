---
title: go-wallet
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.23"

---

# go-wallet

Base URLs:

# Authentication

# Default

## POST 存取款接口

POST /funds/business

> Body 请求参数

```json
{
  "to_account_id": 6,
  "currency": "USD",
  "from_account_id": 6,
  "amount": 1
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Content-Type|header|string| 是 |none|
|body|body|object| 否 |none|
|» owner|body|string| 是 |none|
|» username|body|string| 是 |none|
|» password|body|string| 是 |none|
|» balance|body|integer| 是 |none|
|» currency|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "account": {
    "id": 0,
    "owner": "string",
    "balance": 0,
    "currency": "string",
    "created_at": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» account|object|true|none||none|
|»» id|integer|true|none||none|
|»» owner|string|true|none||none|
|»» balance|integer|true|none||none|
|»» currency|string|true|none||none|
|»» created_at|string|true|none||none|

## POST 创建账户

POST /accounts/creat

> Body 请求参数

```json
{
  "owner": "xhy",
  "currency": "BTC"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST 转账

POST /funds/transfers

> Body 请求参数

```json
{
  "to_account_id": 6,
  "currency": "USD",
  "from_account_id": 6,
  "amount": -1
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "transfer": {
    "id": 0,
    "from_account_id": 0,
    "to_account_id": 0,
    "amount": 0,
    "created_at": "string"
  },
  "from_account": {
    "id": 0,
    "owner": "string",
    "balance": 0,
    "currency": "string",
    "created_at": "string"
  },
  "to_account": {
    "id": 0,
    "owner": "string",
    "balance": 0,
    "currency": "string",
    "created_at": "string"
  },
  "from_entry": {
    "id": 0,
    "account_id": 0,
    "amount": 0,
    "created_at": "string"
  },
  "to_entry": {
    "id": 0,
    "account_id": 0,
    "amount": 0,
    "created_at": "string"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» transfer|object|true|none||none|
|»» id|integer|true|none||none|
|»» from_account_id|integer|true|none||none|
|»» to_account_id|integer|true|none||none|
|»» amount|integer|true|none||none|
|»» created_at|string|true|none||none|
|» from_account|object|true|none||none|
|»» id|integer|true|none||none|
|»» owner|string|true|none||none|
|»» balance|integer|true|none||none|
|»» currency|string|true|none||none|
|»» created_at|string|true|none||none|
|» to_account|object|true|none||none|
|»» id|integer|true|none||none|
|»» owner|string|true|none||none|
|»» balance|integer|true|none||none|
|»» currency|string|true|none||none|
|»» created_at|string|true|none||none|
|» from_entry|object|true|none||none|
|»» id|integer|true|none||none|
|»» account_id|integer|true|none||none|
|»» amount|integer|true|none||none|
|»» created_at|string|true|none||none|
|» to_entry|object|true|none||none|
|»» id|integer|true|none||none|
|»» account_id|integer|true|none||none|
|»» amount|integer|true|none||none|
|»» created_at|string|true|none||none|

## GET 获得账户信息

GET /accounts/{id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "id": 0,
  "owner": "string",
  "balance": 0,
  "currency": "string",
  "created_at": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|integer|true|none||none|
|» owner|string|true|none||none|
|» balance|integer|true|none||none|
|» currency|string|true|none||none|
|» created_at|string|true|none||none|

## GET 获得tx的历史记录

GET /funds/tx

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|account_id|query|integer| 否 |ID 编号|
|page_id|query|integer| 否 |none|
|page_size|query|integer| 否 |none|

> 返回示例

> 200 Response

```json
[
  {
    "id": 0,
    "account_id": 0,
    "amount": 0,
    "created_at": "string"
  }
]
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|integer|true|none||none|
|» account_id|integer|true|none||none|
|» amount|integer|true|none||none|
|» created_at|string|true|none||none|

## GET 获得用户的所有账户

GET /accounts/owners

> Body 请求参数

```yaml
owner: xhy
page_id: 1
page_size: 10

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|owner|query|string| 否 |none|
|page_id|query|string| 否 |none|
|page_size|query|string| 否 |none|
|body|body|object| 否 |none|
|» owner|body|string| 否 |ID 编号|
|» page_id|body|integer| 否 |none|
|» page_size|body|integer| 否 |none|

> 返回示例

> 200 Response

```json
[
  {
    "id": 0,
    "owner": "string",
    "balance": 0,
    "currency": "string",
    "created_at": "string"
  }
]
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|integer|true|none||none|
|» owner|string|true|none||none|
|» balance|integer|true|none||none|
|» currency|string|true|none||none|
|» created_at|string|true|none||none|

# 数据模型


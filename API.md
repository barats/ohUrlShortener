# ohUrlShortener HTTP API

### `/api` 接口权限说明

所有 `/api/*` 接口需要通过 `Bearer Token` 方式验证权限，亦即：每个请求 Header 须携带 

```shell
 Authorization: Bearer {sha256_of_password}
```

`sha256_of_password` 的加密规则，与 `storage/users_storage.go` 中的 `PasswordBase58Hash()` 保持同步

### 1. 新增短链接 `POST /api/url`

接受参数：
1. `dest_url` 目标链接，必填
2. `memo` 备注信息，选填

请求示例：

```shell
curl --request POST \
  --url http://localhost:9092/api/url \
  --header 'Authorization: Bearer EZ2zQjC3fqbkvtggy9p2YaJiLwx1kKPTJxvqVzowtx6t' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data dest_url=http://localhost:9092/admin/dashboard \
  --data memo=dashboard
```

返回结果：

```shell
{
	"code": 200,
	"status": true,
	"message": "success",
	"result": {
		"short_url": "http://localhost:9091/BUUtpbGp"
	},
	"date": "2022-04-10T21:31:29.36559+08:00"
}
```

### 2. 禁用/启用 短链接 `PUT /api/url/:url/change_state`

接受参数：
1. `url` path 参数，指定短链接，必填
2. `enable` 禁用时，传入 false；启用时，传入 true 

请求示例：

```shell
curl --request PUT \
  --url http://localhost:9092/api/url/33R5QUtD/change_state \
  --header 'Authorization: Bearer EZ2zQjC3fqbkvtggy9p2YaJiLwx1kKPTJxvqVzowtx6t' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data enable=false
```

返回结果：

```shell
{
	"code": 200,
	"status": true,
	"message": "success",
	"result": true,
	"date": "2022-04-10T21:31:25.7744402+08:00"
}
```

### 3. 查询短链接统计数据 `GET /api/url/:url`

接受参数：
1. `url` path 参数，指定短链接，必填

请求示例：

```shell
curl --request GET \
  --url http://localhost:9092/api/url/33R5QUtD \
  --header 'Authorization: Bearer EZ2zQjC3fqbkvtggy9p2YaJiLwx1kKPTJxvqVzowtx6t' \
  --header 'Content-Type: application/x-www-form-urlencoded'
```

返回结果：

```shell
{
	"code": 200,
	"status": true,
	"message": "success",
	"result": {
		"short_url": "33R5QUtD",
		"today_count": 3,
		"yesterday_count": 0,
		"last_7_days_count": 0,
		"monthly_count": 3,
		"total_count": 3,
		"d_today_count": 1,
		"d_yesterday_count": 0,
		"d_last_7_days_count": 0,
		"d_monthly_count": 1,
		"d_total_count": 1
	},
	"date": "2022-04-10T21:31:22.059596+08:00"
}
```

### 4. 新建管理员 `POST /api/account`

接受参数：
1. `account` 管理员帐号，必填
2. `password` 管理员密码，必填，最小长度8

请求示例：

```shell
curl --request POST \
  --url http://localhost:9092/api/account \
  --header 'Authorization: Bearer EZ2zQjC3fqbkvtggy9p2YaJiLwx1kKPTJxvqVzowtx6t' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data account=hello1 \
  --data password=12345678
```

返回结果：

```shell
{
	"code": 200,
	"status": true,
	"message": "success",
	"result": null,
	"date": "2022-04-10T21:31:39.7353132+08:00"
}
```

### 5. 修改管理员密码 `PUT /api/account/:account/update`

接受参数：
1. `account` path 参数，管理员帐号，必填
1. `password` 管理员密码，必填，最小长度8

请求示例：

```shell
curl --request PUT \
  --url http://localhost:9092/api/account/hello/update \
  --header 'Authorization: Bearer EZ2zQjC3fqbkvtggy9p2YaJiLwx1kKPTJxvqVzowtx6t' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data password=world123
```

返回结果：

```shell
{
	"code": 200,
	"status": true,
	"message": "success",
	"result": null,
	"date": "2022-04-10T21:31:32.5880538+08:00"
}
```

### 6. 删除短链接 `DELETE /api/url/:url`

接受参数：
1. `url` path 参数，要删除的短链接地址

（此处省略示例）
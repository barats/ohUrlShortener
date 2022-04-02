# ohUrlShortener

一个适合中小型社区网站使用的短链接服务系统，支持短链接成产、查询及302转向，并顺带简单的点击量统计


![Screenshot](screenshot.jpg)

## 配置文件说明
项目根目录下的 `config.ini` 中存放着关于 ohUrlShortener 短链接系统的一些必要配置，请在启动应用之前确保这些配置的正确性

```
[app]

应用是否以 debug 模式启动，主要作用会在go-gin 框架上体现（eg：日志输出等）
debug = false   

短链接系统本地启动端口
port = 9091

短链接系统管理后台本地启动端口
admin_port = 9092

短链接系统的完整 url 前缀，eg：https://t.cn/ 是前缀(不要忘记最后一个/符号)
url_prefix = http://localhost:9091/

[redis] Redis 相关信息
host = 127.0.0.1:56379
database= 0
username=
password=
pool_size = 50

[postgres] 数据库相关信息
host = localhost
port = 55432
user = postgres
password = 0DePm!oG_12Cz^kd_m
database = oh_url_shortener
max_open_conn = 20
max_idle_conn = 5

```

## 所有短域名存储在 Redis 中

所有短链接再系统启动时会以 `Key(short_url) -> Value(original_url)` 的形式存储在 Redis 中。理论上说：如果 Redis 服务器的内存较大的话，存储10w个Key也是可以的。但是，考虑到可扩展性，多封装了一层 `service`，以便需要的时候在业务逻辑层进行自定义扩展，eg：将 key 查询改成数据库查询等。  

## 短域名生产过程相关代码

所在文件 `core/short_url.go` 

```
func GenerateShortLink(initialLink string) (string, error) {
	if utils.EemptyString(initialLink) {
		return "", fmt.Errorf("empty string")
	}
	urlHash, err := utils.Sha256Of(initialLink)
	if err != nil {
		return "", err
	}
	number := new(big.Int).SetBytes(urlHash).Uint64()
	str := utils.Base58Encode([]byte(fmt.Sprintf("%d", number)))
	return str[:8], nil
}
```

## 定时器1分钟清理一次访问日志

所在文件 `main.go` 

```
const ACCESS_LOG_CLEAN_INTERVAL = 1 * time.Minute 

func startTicker() error {
	ticker := time.NewTicker(ACCESS_LOG_CLEAN_INTERVAL)
	for range ticker.C {
		log.Println("[StoreAccessLog] Start.")
		if err := service.StoreAccessLogs(); err != nil {
			log.Printf("Error while trying to store access_log %s", err)
		}
		log.Println("[StoreAccessLog] Finish.")
	}
	return nil
}
```


package conf

/*
===========================================================================================

[mysql.default]
address = "192.168.2.3:3306"
username = "root"
password = "liushuojia"
database = "new_example"

[redis.default]
address = "192.168.2.3:6379"
password = "liushuojia"
db = 30

[token.default]
key = "a1ab2bc3cd4de5ef6fg7g0011223344556677889900123465u11222ser"
#expire = "10m"
expire = "24h"
issuer = "quantify"

[[email.default]]
account = "riskalertbot-1@Balancev.com"
name = "百仁思"
passwd = "rauB7AZWCXv65ztf"
smtp = "smtp.feishu.cn"
port = 465

[minio.default]
address = "192.168.2.3:9000"
access = "BD4SOlUQ8npTbGtLdgPV"
secret = "30KlHODmQxsaSrQxd3wkZSET1UZ4fkct0jz6pEyI"
useSSL = false
bucket = "example"

===========================================================================================
*/
type (
	Token struct {
		Key    string `toml:"key"`
		Expire string `toml:"expire"`
		Issuer string `toml:"issuer"`
	}
	Mysql struct {
		Address  string `toml:"address"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Database string `toml:"database"`
	}
	Redis struct {
		Address  string `toml:"address"`
		Password string `toml:"password"`
		DB       int    `toml:"db"`
	}
	Email struct {
		Account string `toml:"account"`
		Name    string `toml:"name"`
		Passwd  string `toml:"passwd"`
		Smtp    string `toml:"smtp"`
		Port    int    `toml:"port"`
	}
	Minio struct {
		Address string `toml:"address"`
		Access  string `toml:"access"`
		Secret  string `toml:"secret"`
		UseSSL  bool   `toml:"useSSL"`
	}
)

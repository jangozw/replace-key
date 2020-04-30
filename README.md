# 配置项替换工具

根据配置替换配置文件的某些k=>v, 适合不同环境不同的配置文件的部署时候自动修改

# 编译选项

支持跨平台编译



linux 编译生成二进制文件 : ```replace-key```

```bash
make build
```


编译其他平台的见```Makefile``


# 使用demo:

config.ini

```.env
# 本地配置，不会影响其他环境
[database]
schema=mysql
host = abc.com:3306
user   = admin
pwd=1234
dbname=test

[redis]
redis_host=reids:6379
redis_password=
redis_key=
redis_dbNum= 0
cache_expire= 10

# last line must exists

```

replace.json

```.json
{
  "database.host": "test:3309",
  "database.user": "hello99999999999"
}
```

执行替换：
```bash
./replace-key  -source=./file/config.ini -replace=./file/replace.json -output=./file/output.ini
```

输出output.ini:

```..env
# 本地配置，不会影响其他环境

[database]
schema=mysql
host=test:3309
user=hello99999999999
pwd=1234
dbname=test

[redis]
redis_host=reids:6379
redis_password=
redis_key=
redis_dbNum= 0
cache_expire= 10

# last line must exists

```

# 支持任何配置文件

如果是.ini文件, 则匹配section如:[database], 其他文件直接按照关键字替换

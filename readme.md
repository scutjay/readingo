# Readingo

----

Readingo means <b><i>Redis Admin In Golang</i></b>, you could leave the command line alone and use it to
select/update/remove data in redis database via website system. 

### Feature
- Access redis database via web UI
- Support redis and redis cluster
- Adding or removing redis command is allowed and easily
- Two roles provided: readonly and readwrite

### What commands do we support?

We list supported commands in file constant/redis.go. There are two list, one contains read commands and the other one contains write commands.

### Two Authentication ways

When "auth.anonymous" is configured as true, the system could be accessed without login.

When "auth.anonymous" is configured as false, user need to log in with username and password. We provide two roles for user:
- readonly
- readwrite

### Application Architecture

- [gin-gonic/gin](https://github.com/gin-gonic/gin) as web application framework.
- [go-redis/redis/v8](https://github.com/go-redis/redis) as redis client.
- [alecthomas/log4go](https://github.com/alecthomas/log4go) for logging.
- [gopkg.in/yaml.v2]() for reading config.
- [Vue](https://vuejs.org/) for front-end framework.
- [elementUI](https://element.eleme.cn/#/zh-CN) for front-end component library.

### Think aloud

It's an experimental project for me to practice my knowledge about golang and redis. This function is simple and
efficiency may not be good enough, but it didn't mean it could not be use on your project, all opinions are welcomed.

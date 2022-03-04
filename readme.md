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

We list supported commands in file ./constant/redis.go. There are two maps, one for read commands and the other for
write commands, corresponds two role in this system: <i>readonly</i> and <i>readwrite</i>.

### How to configure?

We provide config samples in folder <i>./sample</i>, and the real ones should be put under the runtime directory.

- conf.yml
- log4go.xml

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

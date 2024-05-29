package g

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server struct {
		Mode          string
		Port          string
		DbType        string // mysql | sqlite
		DbAutoMigrate bool
		DbLogMode     string // silent | error | warn | info
		StartTime     string
		MachineID     int64
	}
	Mysql struct {
		Host     string // 服务器地址
		Port     string // 端口
		DSN      string // 高级配置
		Dbname   string // 数据库名
		Username string // 数据库用户名
		Password string // 数据库密码
	}
	Redis struct {
		DB       int    // 指定 Redis 数据库
		Addr     string // 服务器地址:端口
		Password string // 密码
	}
	Log struct {
		DebugFileName string `json:"debugFileName" yaml:"debug_filename"`
		InfoFileName  string `json:"infoFileName" yaml:"info_filename"`
		WarnFileName  string `json:"warnFileName" yaml:"warn_filename"`
		MaxSize       int    `json:"maxsize" yaml:"maxsize"`
		MaxAge        int    `json:"max_age" yaml:"max_age"`
		MaxBackups    int    `json:"max_backups" yaml:"max_backups"`
	}
	Session struct {
		Name   string
		Salt   string
		MaxAge int
	}
	JWT struct {
		Secret string `json:"secret"` // 密钥
		Expire int    `json:"expire"` // 过期时间，单位为 h
		Issuer string `json:"issuer"` // jwt 的发布者
	}
	QiNiu struct {
		Enable    bool   `yaml:"enable"` // 是否用七牛云存储
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
		Bucket    string `json:"bucket"` //存储桶的名字
		CDN       string `json:"cdn"`    //访问图片的地址的前缀
		Zone      string `json:"zone"`   //存储的地区
		Size      string `json:"size"`   //存储的大小限制，单位是MB
	}
	QQ struct {
		AppID    int    `json:"app_id"`
		Key      string `json:"key"`
		Redirect string `json:"redirect"` // 登录之后的回调地址
	}
	Email struct {
		EmailFrom         string `json:"EmailFrom"`
		EmailPassword     string `json:"EmailPassword"`
		SMTPHost          string `json:"SMTPHost"`
		SMTPPort          int    `json:"SMTPPort"`
		EmailWorkerNumber int    `json:"EmailWorkerNumber"`
	}
}

var Conf *Config

func ReadConfig() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		panic("配置文件读取失败: " + err.Error())
	}
	if err := v.Unmarshal(&Conf); err != nil {
		panic("配置文件反序列化失败: " + err.Error())
	}
	log.Println("配置文件内容加载成功.....")
	return Conf
}

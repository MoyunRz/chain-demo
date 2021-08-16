配置文件介绍：
|---config
    |---conf.go     配置文件读取转为实体
    |---conf.toml   配置文件参数设置

配置文件toml格式
        
        [配置名]
        param1 = ..
        param2 = ..

在conf.go文件下添加：
```go
type config struct {
	...
	// 应用配置 
	// App：配置名 app：toml里你设置的配置名称
	App app
}
// app：toml里你设置的配置名称
type app struct {
// toml：toml里你设置的参数名	
Param1  string `toml:"Param1"`
Param1  string `toml:"Param2"`
}

```
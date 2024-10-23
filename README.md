## 须知
运行前，需要在根目录下创建conf文件夹，并新建config.json和plugins.yaml两个配置文件

### config.json
复制后请删掉注释，json不支持注释
```json5
{
  // 这部分配置可参考zeroBot的配置
  "zero": {
    // bot昵称
    "nick_name": ["kohme"],
    // 指令前缀
    "command_prefix": "/",
    // 超级用户
    "super_users": []
  },
  // 正向ws配置
  "ws": {
    "url": "ws://127.0.0.1:3001",
    "token": ""
  },
  // 反向ws配置，不需要启用时留空就好了
  "rws": {
    "url": "ws://127.0.0.1:3002",
    "token": ""
  }

}
```

### plugins.yaml
```yaml
# 插件目录
path: ./plugins 
# 启用的群列表，所有插件的全局配置
groups: [] 
# 各插件配置(下面以core插件示范)
plugins:
  # core插件的配置
  core:
    # 下面说的所有"加载"指的是插件的初始化
    # 插件加载顺序,通过seq的值从小到大依次加载
    # seq相同时，加载顺序不能保证
    seq: 0
    # 是否排除插件(不会加载)
    exclude: false
    # 是否禁用插件(只是禁用功能，但还是会加载)
    disable: false
    # 启用的群列表(为这个插件单独指定启用的群，会屏蔽群设置)
    groups: []
    # conf是对应插件的配置字段，取决于各插件的实现
    conf:
      help_top: 下面是我的所有本领！
      help_tail: 更多本领绝赞学习中,加入github.com/KohmeBot来教会我吧！
    # ... 若有其他键值对，将会作为插件的环境变量传入
```
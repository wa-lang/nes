### 目标: 简化代码, 保持和 WaGo 兼容

- 11月25日, 初始代码 https://github.com/fogleman/nes (5.3k star)
  - 5300 行代码, 本地做了简单测试, 基本操作和声音都是正常的.
- 11月26日, 去掉声音功能(portaudio依赖), 去掉中间状态保存(gob依赖)
  - 4400+ 行代码(减少900行), 无声音/可执行/关机后不保存状态
  - 已经去掉 chan/map 依赖
- 11月26日, 去掉游戏选择视图
  - 4000+ 行代码(减少300行), 启动参数指定游戏
- 11月28日, 去掉 fmt/log/math 依赖
  - 3660+ 行代码(减少300行)
- TODO
  - 去掉 文件系统 依赖
  - 去掉 gl 依赖(最后)
  - 支持 `import "fmt"`
  - 支持 `import "encoding/binary"`
  - 支持 `import "image"`
  - 支持 `import "image/color"`
  - 支持 `import "log"`
  - 支持 `import "os"`


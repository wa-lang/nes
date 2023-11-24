# NES 模拟器

凹语言移植 NES 模拟器. 目标是具备可玩性, 因此优先宣传稳定的版本移植

## Go 语言版 NES

- https://github.com/fogleman/nes (5.3k star), 5000 行代码
  - 本地做了简单测试, 基本操作和声音都是正常的.
- https://github.com/nwidger/nintengo (1.3k star), 6400+3000 行代码
  - 提供了 wasm 版本支持, 基于 Go1.12
- https://github.com/rbaron/awesomenes (270 star), 2500 行代码
  - 在 2019 年移植到了浏览器, 执行时不太稳定.

从代码量/本地测试/Star数量等考虑, 暂定从 `fogleman/nes` 开始移植. 同时也希望能对该实现进行充分测试.

## NES 游戏文件下载

- https://github.com/dream1986/nesrom


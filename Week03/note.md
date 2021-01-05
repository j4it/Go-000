# Summary
* 掌握Goroutine的三个核心点
  * 管控Goroutine的生命周期，即得知道Gorutine什么时候会退出
  * 控制Goroutine的退出（通过close channel/context cancel）
  * 把并行的行为控制扔回给调用者

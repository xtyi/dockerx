docker run command -it
- 创建一个进程，运行 docker init command
- 创建进程的时候使用 clone, 设置 Namespace
- docker init 的代码逻辑会设置 proc 挂载, 然后 exec 函数调用执行 command

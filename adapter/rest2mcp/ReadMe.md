这个package演示如何实现现有rest接口以mcp tool的方式提供服务。

思路参考higress的实现，分离了路由、服务配置的逻辑，仅保留转发。

主要逻辑是：
1. 实现一个mcp server/tool，这个mcp tool和原rest接口是一一对应的关系
2. 当tool被模型请求时，handler中的逻辑会将toolRequest转换成httpRequest请求原rest接口
3. 再把httpResponse转换成toolResponse返回给模型
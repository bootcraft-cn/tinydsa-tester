# TinyDSA Tester

TinyDSA 课程自动测试工具。支持 Java / Python / Go / TypeScript 四种语言的学员解答。

## 方式一：源码构建

```bash
git clone https://github.com/bootcraft-cn/tinydsa-tester
cd tinydsa-tester
go build .
./tinydsa-tester -s dynamic-array -d ~/my-tinydsa
```

**依赖：** Go 1.24+

## 方式二：Docker 镜像

**快速开始**

```bash
cd ~/my-tinydsa  # 你的解答目录
docker pull ghcr.io/bootcraft-cn/tinydsa-tester:latest
docker run --rm --user $(id -u):$(id -g) -v "$(pwd):/workspace" ghcr.io/bootcraft-cn/tinydsa-tester:latest -s dynamic-array -d /workspace
```

**推荐：创建 test.sh 脚本**

在解答目录下创建 `test.sh`：

```bash
#!/bin/bash
docker run --rm --user $(id -u):$(id -g) -v "$(pwd):/workspace" ghcr.io/bootcraft-cn/tinydsa-tester:latest \
  -s "${1:-dynamic-array}" -d /workspace
```

用法：`chmod +x test.sh && ./test.sh stack-and-queue`

## License

MIT

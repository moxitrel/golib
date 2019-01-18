FROM  golang

# 环境变量, $var, ${var}
# ${var:-word}: var未定义，则返回 word
# ${var:+word}: var以定义，则返回 word
#
# 支持环境变量指令：
#   ADD
#   COPY
#   ENV
#   EXPOSE
#   LABEL
#   USER
#   WORKDIR
#   VOLUME
#   STOPSIGNAL
#
#   ONBUILD 支持指令
#ENV  <key>=<val>, ...
#ENV   var   val

# 环境变量, 通过 docker build --build-arg <name> = <value>,
# 若 在FROM之前，只能被FROM使用
#ARG <name>[=<default-value>]
#ARG var=val

RUN ln -s /builds /go/src/github.com

# 复制 本地文件 至 容器; UID, GID为0
# src: must inside the context
#      文件 | 远程文件URL:
#      目录: 仅复制目录中的内容，不复制目录本身
# dst: 创建若不存在
#ADD <src>... <dst>
#ADD ["<src>",... "<dst>"]
#COPY <src>... <dst>                 # src不能为 URL
#COPY ["<src>",... "<dst>"]

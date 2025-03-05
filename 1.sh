$a = $1

if [ $a = "publish" ]; then
    echo "第一个参数是 'publish'，执行相关操作。"
    # 在这里添加需要执行的命令或逻辑
else
    echo "第一个参数不是 'publish'，参数值为：$1"
    # 如果需要，可以在这里处理其他情况
fi
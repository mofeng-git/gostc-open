#!/bin/bash

# 检查是否传入了参数
if [ -z "$1" ]; then
    echo "Usage: $0 <server|gostc> [additional arguments...]"
    exit 1
fi

# 获取脚本参数
TYPE=$1
shift  # 移除第一个参数（server 或 gostc），剩下的参数作为额外参数

# GitHub仓库信息
REPO_OWNER="SianHH"  # 替换为仓库所有者
REPO_NAME="gostc-open"    # 替换为仓库名称

# 目标目录
if [ "$TYPE" = "server" ]; then
    TARGET_DIR="/usr/local/gostc-admin"
    BINARY_NAME="server"
    INSTALL_COMMAND="service install"
elif [ "$TYPE" = "gostc" ]; then
    TARGET_DIR="/usr/local/bin"  # 将 gostc 解压到系统环境
    BINARY_NAME="gostc"
else
    echo "Invalid type. Use 'server' or 'gostc'."
    exit 1
fi

# 创建目标目录（如果不存在）
sudo mkdir -p "$TARGET_DIR"

# 获取当前系统的操作系统类型和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')  # 获取系统类型并转换为小写
ARCH=$(uname -m)                             # 获取系统架构

# 根据架构调整名称
case "$ARCH" in
    "x86_64")
        ARCH="amd64"
        ;;
    "i686"|"i386")
        ARCH="386"
        ;;
    "aarch64"|"arm64")
        ARCH="arm64"
        ;;
    "armv7l"|"armv6l")
        ARCH="arm"
        ;;
    "mips")
        ARCH="mips"
        ;;
    "mips64")
        ARCH="mips64"
        ;;
    "mips64el")
        ARCH="mips64le"
        ;;
    "mipsel")
        ARCH="mipsle"
        ;;
    "riscv64")
        ARCH="riscv64"
        ;;
    "s390x")
        ARCH="s390x"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# 如果是 Windows 系统，调整 OS 名称
if [[ "$OS" == *"mingw"* || "$OS" == *"cygwin"* ]]; then
    OS="windows"
fi

# 获取最新发布的版本信息
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest")

# 提取所有发布文件的URL和名称
ASSETS=$(echo "$LATEST_RELEASE" | grep -oP '"browser_download_url": "\K.*?(?=")')

# 根据类型、系统和架构匹配对应的文件
MATCHED_FILE=""
for ASSET in $ASSETS; do
    if [[ "$ASSET" == *"${TYPE}_${OS}_${ARCH}"* ]]; then
        MATCHED_FILE="$ASSET"
        break
    fi
done

# 下载匹配的文件
if [ -n "$MATCHED_FILE" ]; then
    FILE_NAME=$(basename "$MATCHED_FILE")
    echo "Downloading $FILE_NAME..."
    curl -L -o "$FILE_NAME" "$MATCHED_FILE"

    # 解压文件到目标目录
    echo "Extracting $FILE_NAME to $TARGET_DIR..."
    if [[ "$FILE_NAME" == *.zip ]]; then
        sudo unzip -o "$FILE_NAME" -d "$TARGET_DIR"
    elif [[ "$FILE_NAME" == *.tar.gz ]]; then
        sudo tar -xzf "$FILE_NAME" -C "$TARGET_DIR"
    else
        echo "Unsupported file format: $FILE_NAME"
        exit 1
    fi

    # 修改文件所有者和用户组为 root
    if [ -f "$TARGET_DIR/$BINARY_NAME" ]; then
        sudo chown root:root "$TARGET_DIR/$BINARY_NAME"
        echo "Changed owner and group of $TARGET_DIR/$BINARY_NAME to root."

        # 添加可执行权限
        sudo chmod +x "$TARGET_DIR/$BINARY_NAME"
        echo "Added execute permission to $TARGET_DIR/$BINARY_NAME"
    else
        echo "Binary file $BINARY_NAME not found in $TARGET_DIR."
        exit 1
    fi

    # 如果是 server，运行安装命令
    if [ "$TYPE" = "server" ]; then
        echo "Running installation command for server..."
        sudo "$TARGET_DIR/$BINARY_NAME" $INSTALL_COMMAND "$@"
    fi

    # 清理下载的文件
    rm -f "$FILE_NAME"
    echo "Download, extraction, and installation complete. Files are in $TARGET_DIR."
else
    echo "No matching release file found for ${TYPE}_${OS}_${ARCH}."
fi

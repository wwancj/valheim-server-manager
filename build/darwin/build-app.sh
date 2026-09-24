#!/bin/bash
# 构建 macOS .app 包
# 用法: ./build-app.sh <binary_path> <arch>
# 示例: ./build-app.sh valheim-server-manager-macos-arm64 arm64

set -e

BINARY_PATH=$1
ARCH=$2

if [ -z "$BINARY_PATH" ] || [ -z "$ARCH" ]; then
    echo "Usage: $0 <binary_path> <arch>"
    echo "Example: $0 valheim-server-manager-macos-arm64 arm64"
    exit 1
fi

APP_NAME="ValheimServerManager"
APP_DIR="${APP_NAME}.app"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 清理旧的 .app
rm -rf "${APP_DIR}"

# 创建目录结构
mkdir -p "${APP_DIR}/Contents/MacOS"
mkdir -p "${APP_DIR}/Contents/Resources"

# 复制可执行文件
cp "${BINARY_PATH}" "${APP_DIR}/Contents/MacOS/${APP_NAME}"
chmod +x "${APP_DIR}/Contents/MacOS/${APP_NAME}"

# 复制 Info.plist
cp "${SCRIPT_DIR}/Info.plist" "${APP_DIR}/Contents/Info.plist"

# 创建 PkgInfo
echo -n "APPL????" > "${APP_DIR}/Contents/PkgInfo"

# 复制图标（如果存在）
if [ -f "${SCRIPT_DIR}/AppIcon.icns" ]; then
    cp "${SCRIPT_DIR}/AppIcon.icns" "${APP_DIR}/Contents/Resources/"
    echo "Icon copied."
else
    echo "Warning: AppIcon.icns not found, app will use default icon."
fi

# 打包成 zip
ZIP_NAME="valheim-server-manager-macos-${ARCH}.app.zip"
rm -f "${ZIP_NAME}"
zip -r "${ZIP_NAME}" "${APP_DIR}"

echo ""
echo "=========================================="
echo "Build complete!"
echo "Output: ${ZIP_NAME}"
echo "=========================================="
echo ""
echo "To test:"
echo "  1. Unzip ${ZIP_NAME}"
echo "  2. Double-click ValheimServerManager.app"
echo "  3. Open http://localhost:13256 in browser"

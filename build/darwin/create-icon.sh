#!/bin/bash
# 创建 macOS 应用图标
# 此脚本在 macOS 上运行，生成 Valheim Server Manager 的占位图标

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ICONSET_DIR="${SCRIPT_DIR}/AppIcon.iconset"
OUTPUT_ICNS="${SCRIPT_DIR}/AppIcon.icns"

# 清理
rm -rf "${ICONSET_DIR}"
rm -f "${OUTPUT_ICNS}"

# 创建 iconset 目录
mkdir -p "${ICONSET_DIR}"

# 使用 Python 创建 PNG 图标（带 "V" 字母）
create_icon() {
    local size=$1
    local output=$2

    python3 << EOF
import struct
import zlib

def create_png(width, height, filename):
    # 创建一个带 "V" 字母的圆形图标
    pixels = []

    # 颜色定义 (Valheim 主题色)
    bg_r, bg_g, bg_b = 89, 60, 31  # 深棕色
    fg_r, fg_g, fg_b = 212, 175, 55  # 金色
    border_r, border_g, border_b = 45, 30, 15  # 更深的棕色

    center_x = width / 2
    center_y = height / 2
    radius = min(width, height) / 2 - 2

    for y in range(height):
        row = []
        for x in range(width):
            # 计算到中心的距离
            dx = x - center_x + 0.5
            dy = y - center_y + 0.5
            dist = (dx * dx + dy * dy) ** 0.5

            # 抗锯齿圆形
            if dist > radius + 1:
                row.extend([0, 0, 0, 0])  # 透明
            elif dist > radius - 1:
                # 边缘抗锯齿
                alpha = int(255 * (radius + 1 - dist) / 2)
                alpha = max(0, min(255, alpha))
                row.extend([border_r, border_g, border_b, alpha])
            elif dist > radius - 3:
                # 边框
                row.extend([border_r, border_g, border_b, 255])
            else:
                # 背景色
                row.extend([bg_r, bg_g, bg_b, 255])

                # 绘制 "V" 字母
                # V 的两条斜线
                v_left_start = center_x - width * 0.25
                v_left_end = center_x
                v_right_start = center_x
                v_right_end = center_x + width * 0.25
                v_top = center_y - height * 0.2
                v_bottom = center_y + height * 0.25

                # 计算点到线段的距离
                def point_to_line_dist(px, py, x1, y1, x2, y2):
                    dx = x2 - x1
                    dy = y2 - y1
                    if dx == 0 and dy == 0:
                        return ((px - x1) ** 2 + (py - y1) ** 2) ** 0.5
                    t = max(0, min(1, ((px - x1) * dx + (py - y1) * dy) / (dx * dx + dy * dy)))
                    proj_x = x1 + t * dx
                    proj_y = y1 + t * dy
                    return ((px - proj_x) ** 2 + (py - proj_y) ** 2) ** 0.5

                # 左斜线
                dist_left = point_to_line_dist(x, y, v_left_start, v_top, v_left_end, v_bottom)
                # 右斜线
                dist_right = point_to_line_dist(x, y, v_right_start, v_bottom, v_right_end, v_top)

                line_width = width * 0.06
                if dist_left < line_width or dist_right < line_width:
                    # 在 V 字母上
                    alpha = 255
                    if dist_left < line_width:
                        alpha = int(255 * (line_width - dist_left) / 1.5)
                    elif dist_right < line_width:
                        alpha = int(255 * (line_width - dist_right) / 1.5)
                    alpha = max(0, min(255, alpha))
                    row[-4:] = [fg_r, fg_g, fg_b, alpha]

        pixels.append(bytes(row))

    # 生成 PNG
    def make_chunk(chunk_type, data):
        chunk = chunk_type + data
        return struct.pack('>I', len(data)) + chunk + struct.pack('>I', zlib.crc32(chunk) & 0xffffffff)

    png = b'\x89PNG\r\n\x1a\n'
    png += make_chunk(b'IHDR', struct.pack('>IIBBBBB', width, height, 8, 6, 0, 0, 0))

    raw_data = b''
    for row in pixels:
        raw_data += b'\x00' + row

    png += make_chunk(b'IDAT', zlib.compress(raw_data))
    png += make_chunk(b'IEND', b'')

    with open(filename, 'wb') as f:
        f.write(png)

create_icon(${size}, "${output}")
EOF
}

# 生成不同尺寸的图标
echo "Generating icons..."
SIZES=(16 32 64 128 256 512)
for size in "${SIZES[@]}"; do
    echo "  Creating ${size}x${size}..."
    create_icon $size "${ICONSET_DIR}/icon_${size}x${size}.png"
    # Retina 版本（除了最小的）
    if [ $size -le 256 ]; then
        retina=$((size * 2))
        cp "${ICONSET_DIR}/icon_${size}x${size}.png" "${ICONSET_DIR}/icon_${size}x${size}@2x.png"
    fi
done

# 512@2x = 1024
create_icon 1024 "${ICONSET_DIR}/icon_512x512@2x.png"

# 使用 iconutil 转换为 icns（仅在 macOS 上可用）
if command -v iconutil &> /dev/null; then
    echo "Converting to .icns..."
    iconutil -c icns "${ICONSET_DIR}" -o "${OUTPUT_ICNS}"
    echo "Done: ${OUTPUT_ICNS}"
else
    echo "Warning: iconutil not found (not running on macOS?)"
    echo "Icon files are in: ${ICONSET_DIR}"
    echo "Run this on macOS to create .icns file"
fi

# 清理
rm -rf "${ICONSET_DIR}"

echo ""
echo "Icon creation complete!"

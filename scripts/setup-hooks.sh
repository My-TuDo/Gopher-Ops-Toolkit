#!/bin/sh
#
# 安装 gopt 项目的 Git hooks
# 将 .github/hooks/ 下的 hook 脚本链接到 .git/hooks/
#
# 用法：
#   ./scripts/setup-hooks.sh
#

set -e

HOOKS_SRC="$(cd "$(dirname "$0")/../.github/hooks" && pwd)"
HOOKS_DST="$(cd "$(dirname "$0")/../.git/hooks" && pwd)"

echo "🔗 安装 Git hooks..."
echo "  源目录: $HOOKS_SRC"
echo "  目标目录: $HOOKS_DST"

for hook in "$HOOKS_SRC"/*; do
    name="$(basename "$hook")"
    # 去掉 .sh 后缀，Git hooks 文件名不能带扩展名
    hook_name="${name%.sh}"
    target="$HOOKS_DST/$hook_name"
    ln -sf "$hook" "$target"
    echo "  ✅ $hook_name → $target  (源: $name)"
done

echo ""
echo "✅ Git hooks 安装完成！"
echo ""
echo "  当前启用的 hooks："
for hook in "$HOOKS_SRC"/*; do
    name="$(basename "$hook")"
    echo "    - $name"
done

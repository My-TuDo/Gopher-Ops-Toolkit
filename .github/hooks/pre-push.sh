#!/bin/sh
#
# pre-push hook — 禁止向远程仓库推送，防止密钥/配置泄露
# 安装方式：ln -sf ../../.github/hooks/pre-push.sh .git/hooks/pre-push
#

echo ""
echo "============================================"
echo "  ❌ 推送已阻止：此项目禁止 push 远程仓库"
echo "============================================"
echo ""
echo "  原因：防止 configs/ 等敏感配置泄露"
echo ""
echo "  如需强制推送，请使用："
echo "    git push --no-verify"
echo ""
echo "============================================"
echo ""

exit 1

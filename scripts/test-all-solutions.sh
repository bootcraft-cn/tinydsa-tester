#!/bin/bash
# 批量测试所有 stage 的 solution（Python + Java + Go + TypeScript）
# 用法: ./scripts/test-all-solutions.sh
#
# 分支模型：solution 仓库每种语言一个分支（python / java / go / typescript），
# 脚本通过 git worktree 将各分支 checkout 到临时目录中测试。

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TESTER_DIR="$(dirname "$SCRIPT_DIR")"
SOLUTION_DIR="${TESTER_DIR}/../solution"

# 构建 tester
cd "$TESTER_DIR"
go build -o tinydsa-tester .

# Stage 列表（按课程顺序，slug 需与 tester -s 参数一致）
STAGES=(
    "dynamic-array"         # S01
    "singly-linked-list"    # S02
)

# 语言列表
LANGUAGES=("python" "java" "go" "typescript")

PASSED=0
FAILED=0
SKIPPED=0
TOTAL_TIME=0

echo "=========================================="
echo "  TinyDSA Solution Tester"
echo "=========================================="
echo ""

for lang in "${LANGUAGES[@]}"; do
    echo "--- Language: ${lang} ---"
    echo ""

    # 使用 git worktree 将语言分支 checkout 到临时目录
    # 如果该分支已在主 worktree 中 checkout，则直接使用主目录
    worktree_dir="${SOLUTION_DIR}/.worktree-${lang}"
    use_worktree=true

    current_branch=$(git -C "$SOLUTION_DIR" branch --show-current 2>/dev/null || true)
    if [ "$current_branch" = "$lang" ]; then
        sol_dir="$SOLUTION_DIR"
        use_worktree=false
    else
        if [ -d "$worktree_dir" ]; then
            git -C "$SOLUTION_DIR" worktree remove --force "$worktree_dir" 2>/dev/null || rm -rf "$worktree_dir"
        fi
        if ! git -C "$SOLUTION_DIR" worktree add "$worktree_dir" "$lang" 2>/dev/null; then
            echo "⏭️  [${lang}] SKIPPED - branch not found"
            ((SKIPPED += ${#STAGES[@]}))
            echo ""
            continue
        fi
        sol_dir="$worktree_dir"
    fi

    # TypeScript: 安装依赖（tsx 需要 node_modules）
    if [ "$lang" = "typescript" ] && [ -f "$sol_dir/package.json" ]; then
        printf "📦 [%-24s %10s] Installing deps... " "" "$lang"
        if (cd "$sol_dir" && pnpm install --frozen-lockfile > /dev/null 2>&1); then
            echo "done"
        else
            echo "⚠️  pnpm install failed, trying without --frozen-lockfile"
            (cd "$sol_dir" && pnpm install > /dev/null 2>&1) || true
        fi
    fi

    for stage in "${STAGES[@]}"; do
        printf "🧪 [%-24s %10s] Testing... " "$stage" "$lang"

        start_time=$(python3 -c 'import time; print(time.time())')

        if ./tinydsa-tester -d="$sol_dir" -s="$stage" > /dev/null 2>&1; then
            end_time=$(python3 -c 'import time; print(time.time())')
            elapsed=$(python3 -c "print(f'{$end_time - $start_time:.2f}')")
            echo "✅ PASSED (${elapsed}s)"
            ((PASSED++))
        else
            end_time=$(python3 -c 'import time; print(time.time())')
            elapsed=$(python3 -c "print(f'{$end_time - $start_time:.2f}')")
            echo "❌ FAILED (${elapsed}s)"
            ((FAILED++))
        fi

        TOTAL_TIME=$(python3 -c "print(f'{$TOTAL_TIME + $elapsed:.2f}')")
    done

    # Cleanup worktree（仅清理脚本创建的 worktree）
    if [ "$use_worktree" = true ]; then
        git -C "$SOLUTION_DIR" worktree remove --force "$worktree_dir" 2>/dev/null || true
    fi

    echo ""
done

echo "=========================================="
echo "  Results: $PASSED passed, $FAILED failed, $SKIPPED skipped"
echo "  Total time: ${TOTAL_TIME}s"
echo "=========================================="

if [ $FAILED -gt 0 ]; then
    exit 1
fi

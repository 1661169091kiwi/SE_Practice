@echo off
chcp 65001 >nul
cd /d C:\Users\34517\GolandProjects\SE_Practice

echo 正在添加所有更改...
git add -A

echo 检查是否有未提交的更改...
git diff --cached --quiet
if errorlevel 1 (
    echo 发现未提交的更改，正在提交...
    git commit -m "更新项目内容到目标分支"
)

echo 正在获取远程分支信息...
git fetch origin

echo 正在合并远程目标分支的更改...
git checkout -b temp-merge-branch dev-integrate
git merge origin/初步整合版，包含简单演示界面与粗糙的后端 --no-edit

if errorlevel 1 (
    echo 合并有冲突，正在使用当前分支的内容...
    git checkout --ours .
    git add -A
    git commit -m "合并冲突，保留当前分支内容"
)

echo 正在推送到目标分支...
git push origin temp-merge-branch:refs/heads/初步整合版，包含简单演示界面与粗糙的后端 --force

echo 清理临时分支...
git checkout dev-integrate
git branch -D temp-merge-branch

echo 推送完成！
pause


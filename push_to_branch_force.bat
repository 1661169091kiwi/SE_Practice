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

echo 正在强制推送到目标分支（将覆盖远程分支）...
git push origin dev-integrate:refs/heads/初步整合版，包含简单演示界面与粗糙的后端 --force

if errorlevel 0 (
    echo 推送完成！
) else (
    echo 推送失败，请检查错误信息
)

pause


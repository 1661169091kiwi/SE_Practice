# 设置控制台编码为 UTF-8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

cd C:\Users\34517\GolandProjects\SE_Practice

# 目标分支名
$targetBranch = "初步整合版，包含简单演示界面与粗糙的后端"

# 确保所有更改都已添加
git add -A

# 检查是否有未提交的更改
$status = git status --porcelain
if ($status) {
    Write-Host "发现未提交的更改，正在提交..."
    git commit -m "更新项目内容到目标分支"
}

# 推送到目标分支
Write-Host "正在推送到分支: $targetBranch"
git push origin "dev-integrate:refs/heads/$targetBranch"

Write-Host "推送完成！"


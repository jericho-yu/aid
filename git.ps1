param (
    [string]$operator = $null,
    [string]$arg1 = $null,
    [string]$arg2 = $null
)

# Helper function to run Git commands
function Invoke-GitCommand {
    param (
        [string]$Command
    )
    $gitCommand = "git $Command"
    Invoke-Expression $gitCommand
}

# Main logic
if ($operator -eq "push-tag") {
    $tag = $arg1
    Invoke-GitCommand "tag $tag"
    Invoke-GitCommand "push origin $tag"
}
elseif ($operator -eq "last-tag") {
    Invoke-GitCommand "describe --tags $(git rev-list --tags --max-count=1)"
}
elseif ($operator -eq "help") {
    Write-Output @"
    1、基本用法：. .\git.ps1 执行子程序名称 [参数1, 参数2...]
    2、子程序名称为push-tag  ：创建并推送当前分支内容到tag（参数1表示tag名称）
    3、子程序名称为last-tag  ：获取当前最后一个tag名称
    4、其他子程序名称        ：分为三种情况
        4.1、. .\git.ps1 . 提交说明          ：代表当前分支内容只提交不推送
        4.2、. .\git.ps1 dev 提交说明        ：代表当前内容推送到dev分支
        4.3、. .\git.ps1 dev>master 提交说明 ：代表先推送到dev分支然后合并到master再推送到master
"@
}
else {
    $src_branch, $dst_branch = $operator -split ">"
    $commit = $arg1
    $tag = $arg2

    Invoke-GitCommand "add --all"
    Invoke-GitCommand "commit -m`"$commit`""

    if (-not [string]::IsNullOrEmpty($src_branch)) {
        Invoke-GitCommand "push origin $src_branch"

        if (-not [string]::IsNullOrEmpty($dst_branch)) {
            Invoke-GitCommand "checkout $dst_branch"
            Invoke-GitCommand "merge $src_branch"
            Invoke-GitCommand "push origin $dst_branch"
            Invoke-GitCommand "checkout $src_branch"
        }
    }

    if (-not [string]::IsNullOrEmpty($tag)) {
        Invoke-GitCommand "tag $tag"
        Invoke-GitCommand "push origin $tag"
    }
}

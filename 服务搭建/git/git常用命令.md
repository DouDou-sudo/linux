git init                            #使用当前目录为git仓库，会在当前目录下生成.git目录
git init \<dir\>                    #使用指定目录为git仓库，会在指定目录下生成.git目录

git clone \<remote repo\>           #克隆远端仓库到当前目录，会在当前目录下创建一个目录，目录下包含.git目录，clone仓库不需要init仓库，直接clone
git clone \<remote repo\> \<dir\>   #克隆远端仓库到指定目录

git config -l                       #查看当前的git配置
git config -e                       #编辑git配置，当前仓库
git config -e --global              #针对系统上所有仓库
设置提交代码时的用户信息
git config --global user.name "runoob"
git config --global user.email test@runoob.com
如果去掉 --global 参数只对当前仓库有效

git add .                           #将当前目录下包括子目录下的所有文件更新到暂存区
git add \<file\>                    #只将\<file\>这个文件更新到暂存区
git commit -m ''                    #提交暂存区的本地仓库
git log                             #查看历史提交记录

git rm --cached \<file\>            #从暂存区删除文件，工作区不做出改变，撤销git add 操作
git resert \<commit_id\> \<file\>   #暂存区和本地仓库回退到指定的commit的版本，工作区不会发生改变,撤销指定commit_id版本到目前的git add和git commit操作
git reset HEAD \<file\>             #同上，暂存区和本地仓库回退到上一个版本，工作区不会改变，也可以只指定某一个文件

git checkout .
git checkout -- \<file\>            #用暂存区指定的文件替换当前工作区的文件，这会清除工作区未添加到暂存区的改动
git checkout HEAD .
git checkout HEAD \<file\>          #清空暂存区，用HEAD指向的分支中的全部文件替换暂存区及工作区的文件，这会清除工作区未添加到暂存区的改动，也会清除暂存区未提交的改动

git reset --soft HEAD               #撤销上传的commit记录，仅仅是撤销记录，不会删除工作区的代码
git reset --hard HEAD               #清空暂存区，用HEAD指向的分支中的全部文件替换工作区的文件，这会清除工作区未添加到暂存区的改动，也会清除暂存区未提交的改动

git status                          #查看当前工作区发生的改动，红色为未添加到暂存区的文件，绿色为添加到暂存区但未提交的文件
git status -s                       #上面不加-s输出的文件名默认为index，加上-s可以查看具体的文件名，但是不支持中文，中文为字符编码数字显示

git branch                          #列出分支，并显示当前分支
git checkout \<branch_name\>        #切换分支
git checkout -b \<branch_name\>     #创建并切换分支
git branch -d \<branch_name\>       #删除分支
git merge \<branch_name\>           #合并\<branch_name\>到当前分支

git tag -a \<tagname\>  -m "标签"   #给最新一次提交打上tag

git remote add [<options>] \<name\> \<url\>         #添加远程仓库
本地生成ssh公私钥，然后把公钥粘贴到远程仓库的ssh key里面保存
验证是否互信成功执行如下命令验证
$ ssh -T git@github.com
Hi DouDou-sudo! You've successfully authenticated, but GitHub does not provide shell access.
查看当前的远程库
$ git remote -v
origin  git@github.com:DouDou-sudo/linux.git (fetch)
origin  git@github.com:DouDou-sudo/linux.git (push)
git remote rm \<name\>
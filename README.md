# LearnVocabulary

## 项目介绍

学英语，背单词很费时间。背了又忘，忘了又背。为了更高效的背单词，用go和react编写了这个手机端的webapp，用于在手机浏览器里面背单词。

选择一个单词列表，如四级六级词汇表，乱序生成一个学词计划，每天完成当日的学词任务。防止遗忘，复习是关键。所以，还生成了一个按学习日期排序的复习列表，每天可以选择任意的日期来复习。

## 如何开始

1. 准备数据库

先安装mysql数据库，需要5.7以上版本。然后运行以下命令导入数据表。
```
# 新建数据库
create database learn_vocabulary

# 导入数据表，/path为项目所在路径。
use learn_vocabulary
source /path/LearnVocabulary/database/learn_vocabulary.sql
```

2. 编译react

前端使用react开发，需要nodejs环境来编译。先安装nodejs和npm，再运行以下命令编译。
```
cd /path/LearnVocabulary/front

# 安装依赖
npm install

# 编译
npm run build

# 编译成功后，生成/path/LearnVocabulary/front/build目录
```

3. 编译go

先安装go编译环境，再运行以下命令编译。
```
cd /path/LearnVocabulary/rear

# 假设目标运行环境是64位的linux
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o vocabulary *.go

# 编译成功后，生成vocabulary可执行文件。
```

4. 运行

```
vocabulary -port "5000" -session "Paeh9eivEiJuo1Vu"
    -path "/path/LearnVocabulary/front/build" 
    -mysql "user:password@(127.0.0.1:3306)/learn_vocabulary"

# -port 端口
# -session 用于session的加密
# -path 前端编译生成文件夹的位置
# -mysql 数据库mysql的url路径
```

用浏览器打开 http://localhost:5000/index 就可以看到效果了。
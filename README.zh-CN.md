# 私密空间

[English](./README.md) | 简体中文

## 介绍

这是一款文件加密储存管理器,基于golang+electron+vue技术栈,可用于u盘,普通硬盘并且支持多平台操作系统.

## 快速开始

```bash
# clone 项目
git clone https://github.com/546669204/private-space.git

# 安装依赖 
go get
cd electron
npm install

# 调试 
go build && private-space

cd electron
npm run dev
```


## 构建
```bash
# 构建 golang
go build

# 构建 electron
npm run build:darwin
npm run build:linux
npm run build:mas
npm run build:win32

# 复制golang构建文件到 /electron/build/XXXXXX目录

# 运行
```

## 更新日志
2018-08-06 init


# Private Space

English | [简体中文](./README.zh-CN.md)

This is a file encryption storage manager based on the golang + electronic + VUE technology stack, which can be used for USB drives, common hard drives and supports multi-platform operating systems.

## Getting started

```bash
# clone the project
git clone https://github.com/546669204/private-space.git

# install dependency
go get
cd electron
npm install

# develop 
go build && private-space

cd electron
npm run dev
```


## Build
```bash
#build golang
go build

#build electron
npm run build:darwin
npm run build:linux
npm run build:mas
npm run build:win32

#copy golang build file to /electron/build/XXXXXX

#run
```

## Changelog
2018-08-06 init


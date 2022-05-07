# kinship 亲缘关系判断小程序 <!-- omit in toc -->

## TOC <!-- omit in toc -->

- [部署](#部署)
  - [依赖](#依赖)
  - [安装kinship服务器](#安装kinship服务器)
    - [使用预编译的软件包](#使用预编译的软件包)
    - [使用源码编译安装](#使用源码编译安装)
  - [安装使用依赖](#安装使用依赖)
- [使用](#使用)
  - [启动服务](#启动服务)

## 部署

### 依赖

- build
  - [go1.18](https://go.dev/dl/)
  - git
    - [git-lfs](https://git-lfs.github.io/)
- run
  - python3
    - pandas
  - [ngs-bits](https://github.com/imgag/ngs-bits)

### 安装kinship服务器

#### 使用预编译的软件包

```bash
tar avxf kinship.tar.gz
cd kinship
```

*注意:* 如果您的系统不支持预编译的软件包，请使用源码编译安装。

#### 使用源码编译安装

依赖 [go1.18](https://go.dev/dl/) 和 [git-lfs](https://git-lfs.github.io/)

```bash
git clone https://github.com/liserjrqlxue/kinship.git
cd kinship
go build
git lfs checkout
```

*注意:*
代码在一个不重要的地方使用了泛型，因此依赖 `golang` 的 `1.18` 版本，
如果您的系统暂不支持 `1.18` 版本（比如 `alpine linux` ），
请修改源码降级后编译安装。

### 安装使用依赖

示例使用 `conda` 安装 [ngs-bits](https://github.com/imgag/ngs-bits) 和 `python pandas`

```bash
wget -m https://repo.anaconda.com/miniconda/Miniconda3-py39_4.11.0-Linux-x86_64.sh
sh repo.anaconda.com/miniconda/Miniconda3-py39_4.11.0-Linux-x86_64.sh
. ~/.bashrc
conda config --add channels defaults
conda config --add channels bioconda
conda config --add channels conda-forge
conda create -n ss ngs-bits pandas
```

## 使用

### 启动服务

```bash
conda activate ss
cd kinship
./kinship #-port :9091
```

然后访问 [kinship](http://localhost:9091/kinship) 测试和使用

*注意:*
网址中的 `loclhost` 可换成启动服务器的其它可访问ip或网址，默认端口为 `:9091` ，可通过参数 `-port` 修改。

# ddisk-exporter

> 以前项目叫：
>
> prom-smartctl-exporter
>
> [![Docker Automated build](https://img.shields.io/docker/automated/iodeveloper/prom_smartctlexporter.svg)](https://hub.docker.com/r/iodeveloper/prom_smartctlexporter/)
>
> Docker
>
> [![Docker Hub repository](http://dockeri.co/image/iodeveloper/prom_smartctlexporter)](https://registry.hub.docker.com/u/iodeveloper/prom_smartctlexporter/)
>
> `iodeveloper/prom_smartctlexporter`

安装：

```bash
./install.sh
```

本程序依赖`smartctl`工具

* centos:

  ```bash
  yum install smartmontools -y
  ```

* ubuntu:

  ```bash
  apt get smartmontools -y
  ```


在收集过程中，如果遇到磁盘不支持S.M.A.R.T的，会跳过，只收集支持的磁盘的情况


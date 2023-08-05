# slcx

文档: [[中文](README.md) | [English](doc/README-en.md)]

基于tls加密通道的端口转发工具，绕过安全设备流量检测。

## 如何使用？

以下演示将虚拟机22端口代理到宿主机4443端口

**老板端**

```bash
$ slcx :4444 127.0.0.1:4443 --laoban
```

**牛马端**

```bash
$ slcx 192.168.1.10:4444 127.0.0.1:22 --niuma
```

连接本地4443端口成功登入ssh

```bash
$ ssh vulab@127.0.0.1 -p 4443
```

![](doc/img/ssh-login.png)

## 未来计划

[ ] 支持unix sock

[ ] 支持Udp
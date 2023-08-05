# SLCX

Document: [[中文](README.md) | [English](README-en.md)]

A secure Port forwarding tool that uses tls encryption as the communication shell

## How to use?

The following demonstration proxies virtual machine port 22 to host port 4443

**Boss**

```bash
$ slcx :4443 127.0.0.1:4444 --laoban
```

**Employee**

```bash
$ slcx 192.168.1.10:4443 127.0.0.1:22 --niuma
```

**ssh login**

```bash
$ ssh vulab@127.0.0.1 -p 4443
```

## Future

 [ ] Forward tcp unix

 [ ] Support udp
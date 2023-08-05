package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
)

var secret = []byte("zhimakaimen")
var version = "v1.1"

var certPem = []byte(`-----BEGIN CERTIFICATE-----
MIIDETCCAfkCFH4gCl26J7MJJPCjk9HAlHrpPP3bMA0GCSqGSIb3DQEBCwUAMEUx
CzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRl
cm5ldCBXaWRnaXRzIFB0eSBMdGQwHhcNMjMwODA0MTUxODUzWhcNMzMwODAxMTUx
ODUzWjBFMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UE
CgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAsEP90un068ZDJ3m0PxZ/mwaf1+9b/IAYUfUVDxN75Wa/H6Ak
9xFKD3Htqo/DDFhNh62eiEsesq5mAjHSqgCN9EdHpf1yCZr2YtPJDZUAD2rsEdEd
BJl6hLPNguScoOfE1g6PE8WyUWEH0gkBrglOUnouMkc1xl/LuyBkAPZLzAbHKsxk
bgnQpBPZqngxTTpU969KxWrhrco2LtVdraeF7hmrPbotjCoBN44DPQi2W8vQyxDG
SIu9MOqN0Yz9dIc9SIC+ladExRkkWWuYw6InhWEM/fz0vVrCDOYhrcd4jtcTuYMJ
oEwz6M9eDftZjlRKe8ELA07NOhncr5Ux7YK32wIDAQABMA0GCSqGSIb3DQEBCwUA
A4IBAQCo4sTQN+dusVR2wiFGzMsICBStZo+wCNpjUFNVphyWci2RR9H9I8pjGiGJ
zIvFJh0DziDRtFpx5Qx78PrLdHRcRO4qRHqx6A9NFNi+sePRtiQwdk3lhH3BBZX6
RqdNBtCJ2sDQSwgCkm1r2EjrQA3yGB1wUXwuN9wC33LScm9XAErE9VD6uoKxEOur
m9wY7WIFPazoWe2rH5D9wBMQeKVctTWYnpr+E6xGIoGs5hcFBd10mY97SLHUc6Xx
xDRcwYEScEsYJhjdvf1b/aqCkL/JHwg4SIgebg84wKbwZr/gnkJHBMDL+oev1MqE
SzxDuC6lSPmQ2UoKJ+UUdNeDwa33
-----END CERTIFICATE-----`)

var keyPem = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAsEP90un068ZDJ3m0PxZ/mwaf1+9b/IAYUfUVDxN75Wa/H6Ak
9xFKD3Htqo/DDFhNh62eiEsesq5mAjHSqgCN9EdHpf1yCZr2YtPJDZUAD2rsEdEd
BJl6hLPNguScoOfE1g6PE8WyUWEH0gkBrglOUnouMkc1xl/LuyBkAPZLzAbHKsxk
bgnQpBPZqngxTTpU969KxWrhrco2LtVdraeF7hmrPbotjCoBN44DPQi2W8vQyxDG
SIu9MOqN0Yz9dIc9SIC+ladExRkkWWuYw6InhWEM/fz0vVrCDOYhrcd4jtcTuYMJ
oEwz6M9eDftZjlRKe8ELA07NOhncr5Ux7YK32wIDAQABAoIBAFWdMLQ5HtVGP2f1
HS9nZiOOXCUuBafnL+je+ZORCySpezuJNWsr4qUN2Tp3CW39lerI0a+xjP+M5c87
E++uxURZOtvM9o7jzbjC1WzbamQSswcUYtvREElUbAWftNQSvvCZsbUQ2UPdtUa3
STGXI4u7/KQ/IkdwUCxb1zcAamB2O1FV4T+SSW6oftDeOLUWL686sqek12jdY0hj
n4RxnW71GufN1lZZZW+9AwrUtlLPQ/VmC1n0v5yg63lOOYRxEKhUG15Z9TQTdyng
Nq2goi/9s5fOKZzjwWjJnVT25TH7+SMD1t6eegBDncK5yvZmKL+wF05vw57VgP6L
mE/+voECgYEA5i99tD/nZ2N/jIy9uMjlAw7axx8FDCSuxYHm1IC2Ki8qiT0LcCRY
0di8g0RMzlQ7vYmwWg8qa86VuZC7W3IDzKh0XK2WfJaARPLpHbeG1rFtRrfG3XAB
+cs4bupOdfdcwbweXPr9ga0yh0jEBmdzC4FzqQuavcx5w8iqDLBjmfMCgYEAxAh8
qz8nzrSQlwd9yZ7I9f3Zc/9AlxUA1dUWXJbE80G2XzZ2hgqYBcyPFLjIw78zwqq+
7czn5mEay62qK8zFR+NfJ8orhJQuXhpWufyr0eJt29jznmMMv6nAsWyEs8VXXBlX
FX8cs2asfSlD0y64VYNmOaPAArNbxZ3LPjy+PHkCgYEAr94pFMaeboe+tAv0CHK1
Rpja3rf8i0KbpZuSWVnd1G0k+qYBTzT2B56oqUJNzHbfDfldG/HqX2P/UUAbZcq+
zGcIQKcAC0CR0Bwnr+TAyNfWvek7r4QcflWeUGAEgmz0GbE5IaFCfwltSZ925c20
QF4vmECdgdXwUn3GLtu1vbECgYEAgUmzp2GcyHhQEDuxm0Ls+Ug1nYpQOBRzG5RR
DjmlpWXGeN87pkXgY5q5FwqT5n35Dj98cuu8WjXa7VAZNo/z4Oi1flGZB/SO7CNk
5X8edutmVyOr8N8UrmevnZSbpNDXP/e6doQCEtP5mQaZI/PhQQzLY+J6pY1TTm82
d1KqiMkCgYB9zfEJ9cAW0H7wa9Dj56nJ5cn9owmFDgzM+zA17y9EsaaXRALQuYWj
FCp8//UEzGMUf8Grvuf8EPBD33xyZ3kTBPxtRYHk7RmFSa/dXmiR8fn8/RUCMCjf
bvhgYPeilIO7+5DZ3nPWuAEKQlV3ieGeAcPZPuVvCHcej1eiIZEAig==
-----END RSA PRIVATE KEY-----`)

func config() *tls.Config {
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}
}

func forward(local net.Conn, target net.Conn) {
	defer local.Close()
	defer target.Close()

	go func() {
		log.Println("start forward (local => target)", local.LocalAddr(), "=>", local.RemoteAddr())
		_, err := io.Copy(target, local)
		if err != nil {
			log.Println("error while copying data to remote server:", err)
		}
	}()
	log.Println("start forward (local <= target)", local.LocalAddr(), "<=", local.RemoteAddr())
	_, err := io.Copy(local, target)
	if err != nil {
		log.Println("error while copying data from remote server:", err)
	}
}

func laoban(laobanAddr string, forwardAddr string) {

	l1, err := tls.Listen("tcp", laobanAddr, config())
	if err != nil {
		log.Fatal(err)
	}

	l2, err := net.Listen("tcp", forwardAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer l1.Close()
	defer l2.Close()
	log.Println("laoban listening on", laobanAddr)
	mutex := sync.RWMutex{}
	for {
		mutex.Lock()
		niuma, err := l1.Accept()
		log.Println("niuma report for", niuma.RemoteAddr())
		if err != nil {
			log.Println(err)
		}

		fw, err := l2.Accept()
		log.Println("start forward data for", fw.RemoteAddr())
		if err != nil {
			log.Println(err)
		}

		niuma.Write(secret)

		go forward(niuma, fw)
		mutex.Unlock()
	}
}

func niuma(laobanAddr, forwardAddr string) {
	mutex := sync.RWMutex{}
	for {
		mutex.Lock()
		laoban, err := tls.Dial("tcp", laobanAddr, config())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("start reporting to the laoban")

		fw, err := net.Dial("tcp", forwardAddr)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("start forwarding data")
		for {
			_secret := make([]byte, len(secret))
			_, err := laoban.Read(_secret)

			if err != nil || string(_secret) == string(secret) {
				break
			}
		}

		go forward(fw, laoban)
		mutex.Unlock()
	}
}

func usage() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal("get file path err: ", err)
	}

	fileName := filepath.Base(executable)
	fmt.Println(fmt.Sprintf("slcx %s open source: https://github.com/sechelper/slcx", version))
	fmt.Println("options:")
	fmt.Printf("\t%s laoban_ip:laoban_port forword_ip:forword_port --laoban\n", fileName)
	fmt.Printf("\t%s niuma_ip:niuma_port laoban_ip:laoban_port --niuma\n", fileName)
	fmt.Println("e.g.:")
	fmt.Printf("\t%s :4444 127.0.0.1:4443 --laoban\n", fileName)
	fmt.Printf("\t%s 127.0.0.1:22 192.168.1.10:4444 --niuma\n", fileName)
}

func main() {

	if len(os.Args) < 4 {
		usage()
		return
	}
	switch os.Args[3] {
	case "--laoban":
		laoban(os.Args[1], os.Args[2])
	case "--niuma":
		niuma(os.Args[2], os.Args[1])
	default:
		usage()
	}
}

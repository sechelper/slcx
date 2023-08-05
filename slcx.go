package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var secret = []byte("zhimakaimen")
var version = "v1.1"

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for i := 0; i < length; i++ {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

func generateCertAndKey() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(3650 * 24 * time.Hour)

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: randomString(12) + ".com"},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	certPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	})

	return certPem, keyPem, nil
}

func config() *tls.Config {
	certPem, keyPem, err := generateCertAndKey()
	if err != nil {
		log.Fatal(err)
	}

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
	fmt.Printf("\t%s :4443 127.0.0.1:4444 --laoban\n", fileName)
	fmt.Printf("\t%s 127.0.0.1:22 192.168.1.10:4443 --niuma\n", fileName)
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

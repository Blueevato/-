package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gochat2/common/message"
	"io"
	"net"
)

// 处理数据结构体
type Transfer struct {
	Conn net.Conn
	Buf  [4096]byte
}

func (this *Transfer) ReadPkg() (mess message.Message, err error) {
	//this.Buf := make([]byte, 4096)
	fmt.Println("读取客户端数据...")
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		if err == io.EOF {
			return mess, fmt.Errorf("连接已关闭")
		}
		//	fmt.Println("read len err", err)
		return
	}
	fmt.Println("receive len=", this.Buf[:4])

	//将len转成uint32
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4]) //字节序列解码成 uint32

	//依据pkgLen读取消息内容
	len, err := this.Conn.Read(this.Buf[:pkgLen])
	if uint32(len) != pkgLen || err != nil {
		fmt.Println("read buf err", err)
		return
	}

	fmt.Printf("接收到的数据: %s\n", string(this.Buf[:pkgLen]))
	//把pkgLen反序列化写入Message结构体
	err = json.Unmarshal(this.Buf[:pkgLen], &mess)
	if err != nil {
		fmt.Println("unmarshal err", err)
		return
	}
	return
}

func (this *Transfer) WritePkg(mesdata []byte) (err error) {

	//先发给客户端，先发长度 验证
	var pkgLen uint32
	pkgLen = uint32(len(mesdata))
	//var buf [4]byte                             //uint32 在内存中占用 4字节,32位 = 4字节 × 8位/字节,最大可表示长度：2³² - 1 = 4,294,967,295 字节 ≈ 4GB
	binary.BigEndian.PutUint32(this.Buf[:4], pkgLen) //  编码为大端序（网络字节序）, uint32 编码成字节序列

	n, err := this.Conn.Write(this.Buf[:4]) // 3. 发送长度,bytes是数组，bytes[:4]才是切片
	if err != nil || n != 4 {
		fmt.Println("send len err", err)
		return
	}
	fmt.Printf("成功发送长度: %d 字节\n", len(mesdata))

	// 发送消息本身
	fmt.Printf("内容为: %s\n", mesdata)
	n, err = this.Conn.Write(mesdata)
	if err != nil || n != int(pkgLen) {
		fmt.Println("send len err", err)
		return
	}
	return
}

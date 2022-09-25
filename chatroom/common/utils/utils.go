package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"redis5/chatroom/common/message"
)

//将方法关联到结构体中
type Transfer struct {
	Conn net.Conn
	Buf  [4 * 1024]byte //传输时使用的缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("读取用户发送的数据...")
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		return
	}
	//将buf转为uint32类型
	pkgLen := binary.BigEndian.Uint32(this.Buf[0:4])
	//根据pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}
	//将buf反序列化为message.Message
	err = json.Unmarshal(this.Buf[:pkgLen], &mes) //!!!&mes
	if err != nil {
		fmt.Printf("json.unmarshal err: %v\n", err)
		return
	}
	return
}
func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送长度给对方比较是否一致
	pkgLen := uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[:4], pkgLen)
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Printf("server len c.W err: %v\n", err)
		return
	}
	//发送消息本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Printf("server data c.Werr: %v\n", err)
		return
	}
	return
}

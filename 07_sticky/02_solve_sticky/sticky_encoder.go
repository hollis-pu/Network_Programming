package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-08 10:36
 */

// CusEncoder 定义编码器
type CusEncoder struct {
	w io.Writer
}

// NewEncoder 创建编码器函数
func NewEncoder(w io.Writer) *CusEncoder {
	return &CusEncoder{
		w: w,
	}
}

// Encode 编码，将编码的结果写入到w(io.Writer)
func (encoder CusEncoder) Encode(message string) error {
	// 1.获取message的长度
	l := int32(len(message))

	buf := new(bytes.Buffer) // 构建一个数据包缓冲区

	// 2.在数据包中写入长度
	// 需要二进制的写入操作，需要将数据以bit的形式写入
	if err := binary.Write(buf, binary.LittleEndian, l); err != nil {
		return err
	}

	// 3.将数据主体body写入
	// 方式1：
	// if err := binary.Write(buf, binary.LittleEndian, []byte(message)); err != nil {
	// 	return err
	// }
	// 方式2：
	if _, err := buf.Write([]byte(message)); err != nil {
		return err
	}

	// 4. 利用io.Writer发送数据
	if _, err := encoder.w.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// CusDecoder 定义解码器
type CusDecoder struct {
	r io.Reader
}

// NewDecoder 创建Decoder解码器
func NewDecoder(r io.Reader) *CusDecoder {
	return &CusDecoder{
		r: r,
	}
}

// Decode 解码
func (decoder *CusDecoder) Decode(message *string) error {
	// 1.读取前4个字节
	header := make([]byte, 4)
	headerLen, err := decoder.r.Read(header)
	if err != nil {
		return err
	}
	if headerLen != 4 {
		return errors.New("header is not enough")
	}

	// 2.将前4个字节转换为int32类型，确定了body的长度
	var l int32
	headerBuf := bytes.NewBuffer(header)
	if err := binary.Read(headerBuf, binary.LittleEndian, &l); err != nil {
		return err
	}

	// 3.读取body
	body := make([]byte, l)
	bodyLen, err := decoder.r.Read(body)
	if err != nil {
		return err
	}
	if int32(bodyLen) != l {
		return errors.New("body is not enough")
	}

	// 4.设置message
	*message = string(body)
	return nil
}

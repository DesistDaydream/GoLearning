package main

import (
	"fmt"
)

// KingstonDisk 金士顿牌移动硬盘
type KingstonDisk struct {
	Name string
	Type string
	Data string
}

// NewKingstonDisk is
func NewKingstonDisk() *KingstonDisk {
	return &KingstonDisk{
		Name: "A1",
		Type: "SSD",
		Data: "KingstonDisk fastest SSD",
	}
}

// Start is
func (k *KingstonDisk) Start() {
	fmt.Println("金士顿SSD硬盘已连接")
}

// End is
func (k *KingstonDisk) End() {
	//
}

// OperatorRead is
func OperatorRead(u *KingstonDisk) {
	u.Start()
	fmt.Printf("当前连接的设备信息为：%v\n", u)
	u.End()
}

func main() {
	//
	k := NewKingstonDisk()
	OperatorRead(k)
}

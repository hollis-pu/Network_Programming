package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

/**
* Description:
* @Author Hollis
* @Create 2023-11-02 22:35
 */

type Pool interface {
	Get() (net.Conn, error) //获取连接
	Put(net.Conn) error     //放回连接
	Release() error         //释放池(全部连接)
	Len() int               //有效连接个数
}

type ConnFactory interface {
	Factory(addr string) (net.Conn, error) // 构造连接
	Close(net.Conn) error                  // 关闭连接
	Ping(net.Conn) error                   // 检查连接是否是有效的
}

type PoolConfig struct {
	InitConnNum int           // 初始化连接数
	MaxConnNum  int           // 最大连接数
	MaxIdleNum  int           // 最大空闲连接数
	IdleTimeout time.Duration // 空闲连接超时时间
	Factory     ConnFactory   // 连接工厂
}

// IdleConn 空闲连接结构
type IdleConn struct {
	conn    net.Conn  //连接本身
	putTime time.Time //放回时间
}

type TcpPool struct {
	config         *PoolConfig    // 相关配置
	openingConnNum int            // 开放使用的连接数
	idleList       chan *IdleConn // 空闲的连接队列
	addr           string         // 连接地址
	mut            sync.RWMutex   // 并发安全锁
}

type TcpConnFactory struct{}

func (*TcpConnFactory) Factory(addr string) (net.Conn, error) {
	// 校验参数的合理性
	if addr == "" {
		return nil, errors.New("addr is empty")
	}
	// 建立连接
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	fmt.Printf("client send a conn to %s\n", conn.RemoteAddr())
	// 返回连接对象
	return conn, nil
}
func (*TcpConnFactory) Close(conn net.Conn) error {
	return conn.Close()
}
func (*TcpConnFactory) Ping(net.Conn) error {
	return nil
}

func (pool *TcpPool) Get() (net.Conn, error) {
	// 1.加锁
	pool.mut.Lock()
	defer pool.mut.Unlock()

	// 2.获取空闲连接，若没有则创建连接
	for {
		select {
		// 获取空闲连接
		case idleConn, ok := <-pool.idleList:
			// 判断channel是否被关闭
			if !ok {
				return nil, errors.New("idle list closed")
			}
			// 判断连接是否超时（需要比较pool.config.IdleTimeout和idleConn.putTime）
			if pool.config.IdleTimeout > 0 { // 设置了超时时间（默认值0表示未设置超时时间）
				// 当前时间减去idleConn.putTime是否大于了超时时间（pool.config.IdleTimeout）
				if time.Now().Sub(idleConn.putTime) > pool.config.IdleTimeout {
					// 关闭连接，继续查找下一个连接
					_ = pool.config.Factory.Close(idleConn.conn)
					continue
				}
			}
			// 判断连接是否可用
			if err := pool.config.Factory.Ping(idleConn.conn); err != nil {
				// ping失败，连接不可用 ---> 关闭连接，继续查找
				_ = pool.config.Factory.Close(idleConn.conn)
				continue
			}
			// 找到了可用的空闲连接 ---> 返回连接
			log.Println("get from idle")
			pool.openingConnNum++
			return idleConn.conn, nil
		// 创建连接
		default:
			// 判断是否还可以继续创建（基于开放的连接数是否已经到达了连接池的最大连接数）
			if pool.openingConnNum >= pool.config.MaxConnNum {
				return nil, errors.New("max opening connection")
				// 另一种方案：
				//continue
			}
			// 创建连接
			conn, err := pool.config.Factory.Factory(pool.addr)
			if err != nil {
				return nil, err
			}

			// 正确地创建了可用连接
			log.Println("get from factory")
			pool.openingConnNum++ // 连接计数
			return conn, nil      // 返回连接
		}
	}
}
func (pool *TcpPool) Put(conn net.Conn) error {
	// 1.加锁
	pool.mut.Lock()
	defer pool.mut.Unlock()

	// 2.校验
	if conn == nil { // 传入的连接是否存在
		return errors.New("conn is not exists")
	}
	if pool.idleList == nil { // 空闲连接列表是否存在（如果连接池已关闭，则也将新放入空闲连接列表的连接关闭）
		// 关闭连接
		_ = pool.config.Factory.Close(conn)
		return errors.New("idle list is not exists")
	}

	// 3.放回连接
	select {
	case pool.idleList <- &IdleConn{ // 放回连接
		conn:    conn,
		putTime: time.Now(),
	}:
		// 更新开放的连接数量
		pool.openingConnNum--
		return nil // 只要可以发送成功，任务完成

	default: // 空闲连接列表已满，直接关闭该连接
		_ = pool.config.Factory.Close(conn)
		return nil
	}
}
func (pool *TcpPool) Release() error {
	log.Println("release all conn")

	// 1.加锁（防止其他协程继续操作连接池）
	pool.mut.Lock()
	defer pool.mut.Unlock()

	// 2.确定连接池是否被释放
	if pool.idleList == nil {
		return nil
	}

	// 3.关闭IdleList
	close(pool.idleList)

	// 4.释放全部空闲连接
	// 可以继续接收已关闭channel中的元素
	for idleConn := range pool.idleList {
		_ = pool.config.Factory.Close(idleConn.conn)
	}
	return nil
}
func (pool *TcpPool) Len() int {
	return len(pool.idleList)
}

const (
	defaultInitConnNum = 1
	defaultMaxConnNum  = 10
)

func NewTcpPool(addr string, poolConfig PoolConfig) (*TcpPool, error) {
	// 1.校验参数
	if addr == "" {
		return nil, errors.New("addr is empty")
	}

	if poolConfig.MaxConnNum == 0 { // 合理化最大连接数
		poolConfig.MaxConnNum = defaultMaxConnNum
	}
	if poolConfig.InitConnNum == 0 { // 合理化初始化连接数
		poolConfig.InitConnNum = defaultInitConnNum
	} else if poolConfig.InitConnNum > poolConfig.MaxConnNum {
		poolConfig.InitConnNum = poolConfig.MaxConnNum
	}
	if poolConfig.MaxIdleNum == 0 { // 合理化最大空闲连接数
		poolConfig.MaxIdleNum = poolConfig.InitConnNum
	} else if poolConfig.MaxIdleNum > poolConfig.MaxConnNum {
		poolConfig.MaxIdleNum = poolConfig.MaxConnNum
	}
	if poolConfig.Factory == nil {
		return nil, errors.New("factory is not exists")
	}

	// 2.初始化TcpPool对象
	pool := TcpPool{
		config:         &poolConfig,
		openingConnNum: 0,
		idleList:       make(chan *IdleConn, poolConfig.MaxIdleNum),
		addr:           addr,
		mut:            sync.RWMutex{},
	}

	// 3.初始化连接
	for i := 0; i < poolConfig.InitConnNum; i++ { // 根据InitConnNum来初始化连接的个数
		conn, err := pool.config.Factory.Factory(addr)
		if err != nil { // 初始化连接池失败
			log.Println(err)
			pool.Release() // 释放可能存在的连接
			return nil, err
		}

		pool.idleList <- &IdleConn{ // 连接创建成功，加入到空闲连接队列中
			conn:    conn,
			putTime: time.Now(),
		}
	}

	// 4.返回
	return &pool, nil
}

func TCPServerPool() {
	log.Println("服务器端开始监听8888端口...")

	host := "localhost"
	port := 8888

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Println("listener err=", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go HandleConnPool(conn)
	}
}

func HandleConnPool(conn net.Conn) {
	log.Printf("accept conn from %s\n", conn.RemoteAddr())
	defer func() {
		log.Println("conn be closed")
		conn.Close()
	}()
}

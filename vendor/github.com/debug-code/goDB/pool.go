package goDB

import (
	"errors"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
)

var (
	ErrInvalidConfig = errors.New("invalid pool config")
	ErrPoolClosed    = errors.New("pool closed")
)

type factory func() (*gorm.DB, error)
type close func(*gorm.DB) error

type Pool interface {
	Acquire() (*gorm.DB, error) //获取资源
	Release(*gorm.DB) error     //释放资源
	Close(*gorm.DB) error       //关闭资源
	ShutDown() error            //关闭池
}

type GenericPool struct {
	sync.RWMutex
	pool        chan *gorm.DB
	maxOpen     int
	numOpen     int
	minOpen     int           //池中最少资源
	closed      bool          //池关闭状态
	maxLifeTime time.Duration //生存时间
	factory     factory       //create function
	close       close         //关闭连接使用
}

func NewGenericPool(minOpen, maxOpen int, maxLifeTime time.Duration, factory factory, close close) (*GenericPool, error) {

	if maxOpen <= 0 || minOpen > maxOpen {
		return nil, ErrInvalidConfig
	}

	p := &GenericPool{
		maxOpen:     maxOpen,
		minOpen:     minOpen,
		maxLifeTime: maxLifeTime,
		factory:     factory,
		pool:        make(chan *gorm.DB, maxOpen),
		close:       close,
	}

	for i := 0; i < minOpen; i++ {
		closer, err := factory()
		if err != nil {
			continue
		}
		p.numOpen++
		p.pool <- closer
	}
	return p, nil
}

func (p *GenericPool) GetLife() int {
	//p.Lock()
	min := p.numOpen
	//p.Unlock()
	return min
}

func (p *GenericPool) Acquire() (*gorm.DB, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}

	for {
		closer, err := p.getOrCreate()
		if err != nil {
			return nil, err
		}

		// todo maxLifeTime 处理

		return closer, nil
	}
}

func (p *GenericPool) getOrCreate() (*gorm.DB, error) {
	select {
	case closer := <-p.pool:
		return closer, nil
	default:
	}
	p.Lock()
	defer p.Unlock()
	if p.numOpen >= p.maxOpen {
		closer := <-p.pool
		//p.Unlock()
		return closer, nil
	}

	// new connect
	closer, err := p.factory()
	if err != nil {
		p.Unlock()
		return nil, err
	}
	p.numOpen++
	//p.Unlock()
	return closer, nil

}

// 释放单个资源到连接池
func (p *GenericPool) Release(closer *gorm.DB) error {
	if p.closed {
		return ErrPoolClosed
	}

	//if p.numOpen > p.maxOpen {
	//	p.Lock()
	//	err := p.close(closer)
	//	if err != nil {
	//		return err
	//	}
	//	p.numOpen--
	//	p.Unlock()
	//	return nil
	//} else {
	//
	//}
	p.pool <- closer
	return nil
}

// 关闭单个资源
func (p *GenericPool) Close(closer *gorm.DB) error {
	p.Lock()
	//closer.(io.Closer).Close()
	p.close(closer)
	p.numOpen--
	p.Unlock()
	return nil
}

// 关闭连接池，释放所有资源
func (p *GenericPool) Shutdown() error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	//close(p.pool)
	for closer := range p.pool {
		//closer.(io.Closer).Close()
		p.close(closer)
		p.numOpen--
	}
	p.closed = true
	p.Unlock()
	return nil
}

package balance

import "sync/atomic"

type LoadBalancer interface {
	Address() string
}

type loadBalancer struct {
	addresses []string
	counter   atomic.Int32
}

func (l *loadBalancer) Address() string {
	val := l.counter.Add(1) - 1
	return l.addresses[val%int32(len(l.addresses))]
}

func NewSimpleLoadBalancer(addresses ...string) LoadBalancer {
	return &loadBalancer{
		addresses: addresses,
	}
}

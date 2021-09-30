package schedule

import "strings"

type Resolver interface {
	Resolve(runner string) (schema string, addrs []string, err error)
}

func NewDirectResolver() Resolver {
	return DirectResolver{}
}

type DirectResolver struct {
}

func (r DirectResolver) Resolve(runner string) (schema string, addrs []string, err error) {
	array := strings.SplitN(runner, "://", 2)
	schema, addrs = array[0], []string{runner}
	return
}

type NacosResolver struct {
}

func (r NacosResolver) Resolve(runner string) (schema string, addrs []string, err error) {
	panic("implement me")
}

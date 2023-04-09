package ff1

import (
	"errors"
	"sync"
)

var (
	_self *transform
	once  sync.Once
)

type transform struct {
	intToChar map[int]string
	charToInt map[string]int
	sync.Once
}

func GetTransform() *transform {
	if _self == nil {
		once.Do(func() {
			newTransform()
		})
	}
	return _self
}

func newTransform() *transform {
	_self = &transform{}
	hashmapSize := 10 + 26*2 + 0x9fa5 - 0x4e00 + 2
	_self.charToInt = make(map[string]int, hashmapSize)
	_self.intToChar = make(map[int]string, hashmapSize)
	i := 0
	for c := '0'; c <= '9'; c++ {
		_self.charToInt[string(c)] = i
		_self.intToChar[i] = string(c)
		i += 1
	}
	for c := 'a'; c <= 'z'; c++ {
		_self.charToInt[string(c)] = i
		_self.intToChar[i] = string(c)
		i += 1
	}
	for c := 'A'; c <= 'Z'; c++ {
		_self.charToInt[string(c)] = i
		_self.intToChar[i] = string(c)
		i += 1
	}
	for k := 0x4e00; k <= 0x9fa5; k++ {
		_self.charToInt[string(k)] = i
		_self.intToChar[i] = string(k)
		i += 1
	}

	//space 0020
	space := string(0x0020)
	_self.charToInt[space] = i
	_self.intToChar[i] = space
	i += 1

	return _self
}
func (t *transform) TransformInt(data []int) (string, error) {

	var result string
	for i := 0; i < len(data); i++ {
		v, ok := t.intToChar[data[i]]
		if !ok {
			return "", errors.New("invalid input data")
		}
		result += v
	}
	return result, nil
}
func (t *transform) TransformString(data string) ([]int, error) {

	var result []int
	for _, c := range data {
		v, ok := t.charToInt[string(c)]
		if !ok {
			return nil, errors.New("invalid input data")
		}
		result = append(result, v)
	}
	return result, nil
}

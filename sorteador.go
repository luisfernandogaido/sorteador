package sorteador

import (
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

type Sorteio struct {
	min       int
	max       int
	sorteados map[int]bool
}

func NewSorteio(min, max int) (Sorteio, error) {
	if min < 0 {
		return Sorteio{}, errors.New("min deve ser positivo")
	}
	if max < 0 {
		return Sorteio{}, errors.New("max deve ser positivo")
	}
	if min > max {
		return Sorteio{}, errors.New("min não pode ser maior que max")
	}
	return Sorteio{min: min, max: max, sorteados: make(map[int]bool)}, nil
}

func (s Sorteio) Proximo() (int, error) {
	if len(s.sorteados) == (s.max - s.min + 1) {
		return 0, errors.New("fim do sorteio")
	}
	for {
		numero := s.min + rand.Intn(s.max+1)
		if !s.sorteados[numero] {
			s.sorteados[numero] = true
			return numero, nil
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

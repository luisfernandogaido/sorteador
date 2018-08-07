package app

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const NomeSessao = "GO_SESSION"

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	muSessoes = sync.Mutex{}
	rd        redis.Conn
)

type Sessao struct {
	chave   string
	Valores map[string]interface{}
}

func SessaoRedis(conn redis.Conn) {
	rd = conn
}

func (s Sessao) Salva() error {
	bytes, err := json.Marshal(s.Valores)
	if err != nil {
		return err
	}
	muSessoes.Lock()
	defer muSessoes.Unlock()
	_, err = rd.Do("SET", NomeSessao+":"+s.chave, bytes)
	return err
}

func (s Sessao) Int(chave string) (int, bool) {
	muSessoes.Lock()
	defer muSessoes.Unlock()
	if valor, ok := s.Valores[chave]; ok {
		return int(valor.(float64)), true
	} else {
		return 0, false
	}
}

func (s Sessao) Float64(chave string) (float64, bool) {
	muSessoes.Lock()
	defer muSessoes.Unlock()
	if valor, ok := s.Valores[chave]; ok {
		return valor.(float64), true
	} else {
		return 0.0, false
	}
}

func SessaoIni(w http.ResponseWriter, r *http.Request) (Sessao, error) {
	sess := Sessao{}
	cookie, err := r.Cookie(NomeSessao)
	if err != nil {
		return novaSessao(w)
	}
	muSessoes.Lock()
	defer muSessoes.Unlock()
	bytes, err := redis.Bytes(rd.Do("GET", NomeSessao+":"+cookie.Value))
	if err != nil {
		return novaSessao(w)
	}
	if err := json.Unmarshal(bytes, &sess.Valores); err != nil {
		return Sessao{}, err
	}
	sess.chave = cookie.Value
	return sess, nil
}

func geraValorSessao() (string, error) {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	return fmt.Sprintf("%x", bytes), err
}

func novaSessao(w http.ResponseWriter) (Sessao, error) {
	sess := Sessao{}
	valor, err := geraValorSessao()
	if err != nil {
		return Sessao{}, err
	}
	cookie := &http.Cookie{
		Name:  NomeSessao,
		Value: valor,
	}
	http.SetCookie(w, cookie)
	sess.chave = valor
	sess.Valores = make(map[string]interface{})
	return sess, nil
}

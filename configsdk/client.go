// Author: Sunshibao <664588619@qq.com>
package configsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

type Decoder func(data []byte, v interface{}) error

type Client interface {
	Use(l Loader) Client
	Watch(key string, p interface{}, decoder Decoder) error
	WatchJSON(key string, p interface{}) error
	WatchYAML(key string, p interface{}) error
	Load(key string) ([]byte, error)
	LoadJSON(key string, v interface{}) error
	LoadYAML(key string, v interface{}) error
}

func New(dir string) (Client, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	if info, err := os.Stat(abs); err == nil {
		if !info.IsDir() {
			return nil, errors.New("config: dir is not a directory")
		}
	} else {
		if os.IsNotExist(err) {
			err = os.MkdirAll(abs, 0766)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	fs, err := newFiles(abs)
	if err != nil {
		return nil, err
	}
	return &client{fs: fs}, nil
}

func MustNew(dir string) Client {
	client, err := New(dir)
	if err != nil {
		panic(err)
	}
	return client
}

type client struct {
	mutex  sync.RWMutex
	fs     *files
	loader []Loader
	items  []*watchable
	cache  map[string][]byte
	once   sync.Once
}

func (c *client) Use(l Loader) Client {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.loader = append(c.loader, l)
	return c
}

func (c *client) Watch(key string, p interface{}, decoder Decoder) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.watch(key, p, decoder)
}

func (c *client) WatchJSON(key string, p interface{}) error {
	return c.Watch(key, p, json.Unmarshal)
}

func (c *client) WatchYAML(key string, p interface{}) error {
	return c.Watch(key, p, yaml.Unmarshal)
}

func (c *client) watch(key string, p interface{}, f func([]byte, interface{}) error) error {
	t := reflect.TypeOf(p)
	if kind := t.Kind(); kind != reflect.Ptr || !reflect.New(t.Elem()).CanInterface() {
		return errors.New(`watcher: invalid config object pointer`)
	}
	// 首次监听需完成配置项加载
	data, err := c.load(key)
	if err != nil {
		return err
	}

	w := &watchable{key, p, t, f, nil}
	if err := f(data, w.p); err != nil {
		return err
	}
	w.cache = copyBytes(data)
	c.items = append(c.items, w)
	c.once.Do(c.doLoop)
	return nil
}

func (c *client) doLoop() {
	go func() {
		ch := make(chan struct{})
		go func() {
			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			defer signal.Stop(c)
			<-c
			close(ch)
		}()
		for {
			if err := c.call(func() error { return c.loop(ch) }); err != nil {
				GetLogger().Errorf("config: watch loop error: %s", err)
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}
	}()
}

func (c *client) call(f func() error) (err error) {
	defer func() {
		if v := recover(); v != nil {
			switch o := v.(type) {
			case string:
				err = fmt.Errorf("config panic: %s", o)
			case error:
				err = fmt.Errorf("config panic: %s", o.Error())
			default:
				err = fmt.Errorf("config panic: %v", v)
			}
		}
	}()
	err = f()
	return
}

func (c *client) loop(ch <-chan struct{}) (err error) {
	d := time.Second * 15
	t := time.NewTimer(d)
	defer t.Stop()

	for {
		select {
		case <-ch:
			return nil
		case <-t.C:
			if err := c.call(c.doWatch); err != nil {
				return err
			}
			t.Reset(d)
		}
	}
}

func (c *client) doWatch() (err error) {
	var items []*watchable
	c.mutex.RLock()
	if len(c.items) > 0 {
		items = c.items[:]
	}
	c.mutex.RUnlock()
	if len(items) > 0 {
		var es errs
		for i, j := 0, len(items); i < j; i++ {
			data, err := c.Load(items[i].key)
			if err != nil {
				es = append(es, fmt.Errorf("load item %s error: %s", items[i].key, err))
				continue
			}
			if bytes.Equal(data, items[i].cache) {
				continue
			}
			o := items[i].new()
			if err := items[i].f(data, o); err != nil {
				es = append(es, fmt.Errorf("decode item %s error: %s", items[i].key, err))
				continue
			}
			items[i].cache = copyBytes(data)
			items[i].swap(o)
			GetLogger().Infof("config: swap item %s done", items[i].key)
		}
		if len(es) > 0 {
			for _, e := range es {
				GetLogger().Errorf("config: watch error: %s", e)
			}
		}
	}
	return
}

func (c *client) Load(key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.load(key)
}

func (c *client) LoadJSON(key string, v interface{}) error {
	data, err := c.Load(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (c *client) LoadYAML(key string, v interface{}) error {
	data, err := c.Load(key)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

func (c *client) load(key string) ([]byte, error) {
	var loaded []byte
	var loadErrs errs
	for i := len(c.loader) - 1; i >= 0; i-- {
		data, err := c.loader[i].Load(key)
		if err != nil {
			if err != ErrNotFound {
				// 加载器失败，记录错误，尝试其他加载器
				loadErrs = append(loadErrs, err)
			}
			continue
		}
		if len(data) == 0 {
			continue
		}
		if err := c.fs.set(key, data); err != nil {
			// 将远程配置存储到本地文件系统中失败时仅报告错误，返回远端最新配置
			loadErrs = append(loadErrs, err)
		}
		loaded = data
		break
	}
	if len(loadErrs) > 0 {
		GetLogger().Errorf("config: Load %s error: %s", key, loadErrs.Error())
	}
	if len(loaded) > 0 {
		return loaded, nil
	}
	// 从本地文件系统查询配置项
	if data, found := c.fs.get(key); found && len(data) > 0 {
		return data, nil
	}
	return nil, ErrNotFound
}

type watchable struct {
	key   string
	p     interface{}
	t     reflect.Type
	f     func([]byte, interface{}) error
	cache []byte
}

func (w *watchable) new() interface{} {
	return reflect.New(w.t.Elem()).Interface()
}

func (w *watchable) swap(o interface{}) {
	reflect.ValueOf(w.p).Elem().Set(reflect.ValueOf(o).Elem())
}

type files struct {
	dir   string
	mutex sync.RWMutex
	items map[string][]byte
}

func newFiles(dir string) (*files, error) {
	r := &files{dir: dir}
	if err := r.sync(); err != nil {
		return nil, err
	}
	return r, nil
}

func (f *files) sync() error {
	matches, err := filepath.Glob(filepath.Join(f.dir, "*"))
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return nil
	}
	items := make(map[string][]byte)
	for _, path := range matches {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			continue
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		items[info.Name()] = data
	}
	if len(items) > 0 {
		f.mutex.Lock()
		f.items = items
		f.mutex.Unlock()
	}
	return nil
}

func (f *files) get(key string) (data []byte, found bool) {
	f.mutex.RLock()
	data, found = f.items[key]
	f.mutex.RUnlock()
	return
}

func (f *files) set(key string, data []byte) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	old, found := f.items[key]
	// 仅在文件不存在或者文件内容不一致时才进行内容缓存和落本地文件系统
	if !found || !bytes.Equal(old, data) {
		err := ioutil.WriteFile(filepath.Join(f.dir, key), data, 0666)
		if err != nil {
			return err
		}
		f.items[key] = data
	}
	return nil
}

type errs []error

func (es errs) Error() string {
	switch n := len(es); n {
	case 0:
		return "<empty errors>"
	case 1:
		return es[0].Error()
	default:
		ss := make([]string, 0, n)
		for i := 0; i < n; i++ {
			ss = append(ss, es[i].Error())
		}
		return strings.Join(ss, "; ")
	}
}

func copyBytes(src []byte) []byte {
	if len(src) == 0 {
		return []byte{}
	}
	dst := make([]byte, len(src))
	copy(src, dst)
	return dst
}

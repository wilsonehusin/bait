package bait

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

type RootConfig struct {
	Config []*BaitConfig `yaml:"config"`
}

type BaitConfig struct {
	Request string `yaml:"request"`
	Workdir string `yaml:"workdir"`
	Command string `yaml:"command"`

	lock sync.Mutex
}

func (b *BaitConfig) Lock() {
	b.lock.Lock()
}

func (b *BaitConfig) Unlock() {
	b.lock.Unlock()
}

func (b *BaitConfig) Cmd(ctx context.Context) *exec.Cmd {
	words := strings.Split(b.Command, " ")
	cmd := exec.CommandContext(ctx, words[0], words[1:]...)
	cmd.Dir = b.Workdir
	return cmd
}

func NewConfigFromFile(str string) (*RootConfig, error) {
	data, err := os.ReadFile(str)
	if err != nil {
		return nil, err
	}
	return NewConfig(data)
}

func NewConfig(b []byte) (*RootConfig, error) {
	c := RootConfig{}
	err := yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

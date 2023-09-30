package bait_test

import (
	"context"
	"testing"

	"go.husin.dev/bait"
)

func TestNewConfig(t *testing.T) {
	src := []byte(`---
config:
  - request: "/hello.world"
    workdir: "/www"
    command: "echo hello.world"
`)
	c, err := bait.NewConfig(src)
	if err != nil {
		t.Fatal(err)
	}

	conf := c.Config
	if conf[0].Request != "/hello.world" {
		t.Fatalf("request: %v != /hello.world", conf[0].Request)
	}
	if conf[0].Workdir != "/www" {
		t.Fatalf("workdir: %v != /www", conf[0].Workdir)
	}
	if conf[0].Command != "echo hello.world" {
		t.Fatalf("command: %v != echo hello.world", conf[0].Command)
	}
}

func TestConfigCommandArgs(t *testing.T) {
	conf := &bait.BaitConfig{
		Request: "lalalala",
		Workdir: "lililili",
		Command: "doesnotexistwoohoo foo bar",
	}
	cmd := conf.Cmd(context.Background())
	if cmd.Path != "doesnotexistwoohoo" {
		t.Fatalf("path: %v != doesnotexistwoohoo", cmd.Path)
	}
	if cmd.Args[0] != "doesnotexistwoohoo" {
		t.Fatalf("args[0]: %v != doesnotexistwoohoo", cmd.Args[0])
	}
	if cmd.Args[1] != "foo" {
		t.Fatal()
	}
	if cmd.Args[2] != "bar" {
		t.Fatal()
	}
}

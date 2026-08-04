package main

import (
	"time"

	. "github.com/apiarytech/royaljelly/core"
)

type HelloWorldProgram struct {
	MyButton BOOL //(* Input from a switch *)
	MyLamp   BOOL //(* Output to an indicator light *)
}

func (p *HelloWorldProgram) Init() {
	p.MyButton = false
	p.MyLamp = false
}

func (p *HelloWorldProgram) Logic(now time.Time) {
	if p.MyButton {
		p.MyLamp = true
	} else {
		p.MyLamp = false
	}
}

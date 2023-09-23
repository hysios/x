package plugin

import (
	"testing"

	"github.com/tj/assert"
)

type Echo interface {
	Echo(s string) string
	isPlugin()
}

type EchoPlugin struct {
}

func (p *EchoPlugin) Echo(s string) string {
	return s
}

// isPlugin
func (*EchoPlugin) isPlugin() {}

func TestHost(t *testing.T) {
	var host = Create[Echo]()
	host.Install(Desc{Name: "echo"}, &EchoPlugin{})

	agent, err := MakeAgent(host)
	assert.NoError(t, err)
	agent.Echo("hello")

}

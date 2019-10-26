/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 16-12-2018
 * |
 * | File Name:     instance.go
 * +===============================================
 */

package models

import (
	"fmt"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/toskatok/lg/core"
	"github.com/toskatok/lg/generators"
)

// Instance represents a load generator instance. Each instance has a runner with
// a generator that generates messages.
// R and G are public and consumers can use their public methods in their advanced use cases.
type Instance struct {
	R core.Runner
	G generators.Generator

	// message is used for populating the data templates.
	message struct {
		Count int64 // send message counter
	}
}

// NewInstance creates new load generator instance with given configuration
func NewInstance(config Config, rate time.Duration, destination string) (*Instance, error) {
	instance := &Instance{}

	// Create parent template with some useful function
	tmpl := template.New("lg").Funcs(template.FuncMap{
		"now":   time.Now,
		"randn": rand.Intn,
	})

	// Read data from the given data templates, and will prase the template if it exists.
	for i, d := range config.Messages {
		for k, v := range d {
			if s, ok := v.(string); ok { // parses all strings
				p, err := tmpl.New(fmt.Sprintf("lg-%d-%s", i, k)).Parse(s)
				if err != nil {
					return nil, err
				}

				d[k] = p
			}
		}
	}

	var err error

	instance.G, err = generators.Get(config.Generator.Name, config.Generator.Info)
	if err != nil {
		return nil, err
	}

	// Runner creation
	instance.R, err = core.NewRunner(core.RunnerConfig{
		Generator: instance.G,
		Duration:  rate,
		Pick: func() interface{} { // runs on each message
			instance.message.Count++

			d := make(map[string]interface{})
			for k, v := range config.Messages[rand.Intn(len(config.Messages))] {
				if tmpl, ok := v.(*template.Template); ok {
					var b strings.Builder
					if err := tmpl.Execute(&b, instance.message); err != nil {
						continue
					}
					d[k] = b.String()
				} else {
					d[k] = v
				}
			}
			return d
		},
		URL:   destination,
		Token: config.Token,
	})
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// Run runs the instance (please note that runners run in new go routine by default)
func (i *Instance) Run() {
	i.R.Run()
}

// Stop stops the instance
func (i *Instance) Stop() {
	i.R.Stop()
}

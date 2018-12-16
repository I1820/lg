/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 20-11-2017
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/I1820/lg/models"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// config variable contains current user configuration
var config models.Config

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // looking for config in the working directory

	if err := viper.ReadInConfig(); err != nil { // find and read the config file
		log.Fatalf("config file: %s \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil { // parse configuration into Config structure
		log.Fatal(err)
	}

	app := &cli.App{
		Name:        "MQTT-LG",
		Description: "MQTT based Load Generator",
		Authors: []cli.Author{
			cli.Author{
				Name:  "Parham Alvani",
				Email: "parham.alvani@gmail.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "destination",
				Value: "mqtt://127.0.0.1:1883",
				Usage: "scheme://(host or host:port)",
			},
			&cli.DurationFlag{
				Name:  "rate",
				Value: 1 * time.Millisecond,
				Usage: "Send interval",
			},
		},
		Action: func(c *cli.Context) error {
			i, err := models.NewInstance(config, c.Duration("rate"), c.String("destination"))
			if err != nil {
				return err
			}

			// print generator information
			fmt.Println(">>> Generator")
			fmt.Printf("%+v\n", i.G)
			fmt.Println(">>>")

			// runs the instance
			i.Run()

			// prints report in 1 second intervals
			go func() {
				for {
					time.Sleep(1 * time.Second)
					fmt.Printf("%s -> %d\n", config.Generator.Name, i.R.Count())
				}
			}()

			// Set up channel on which to send signal notifications.
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, os.Interrupt, os.Kill)

			<-sigc
			i.Stop()
			fmt.Printf("Total packets send to %s: %d\n", config.Generator.Name, i.R.Count())

			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
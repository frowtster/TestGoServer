package proc_EventHandler

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"tnh-simple-server/t_util"

	"github.com/streadway/amqp"
)

type Exchange struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Arguments struct {
	Xdeadletterexchange   string `json:"x-dead-letter-exchange"`
	Xdeadletterroutingkey string `json:"x-dead-letter-routing-key"`
	Xmessagettl           int    `json:"x-message-ttl"`
	Xexpires              int    `json:"x-expires"`
}

type Queue struct {
	Name      string    `json:"name"`
	Arguments Arguments `json:"arguments"`
}

type Option struct {
	Uid        string `json:"uid"`
	AccessTime string `json:"accessTime"`
}

type Bind struct {
	Key    string `json:"key"`
	Option Option `json:"option"`
}

type Rabbit struct {
	Host      string   `json:"host"`
	Port      int      `json:"port"`
	User      string   `json:"user"`
	Password  string   `json:"password"`
	Vhost     string   `json:"vhost"`
	Heartbeat int      `json:"heartbeat"`
	Policies  string   `json:"policies"`
	Exchange  Exchange `json:"exchange"`
	Queue     Queue    `json:"queue"`
	Bind      Bind     `json:"bind"`
}

type ConfigEventHandler struct {
	t_util.ConfigInfo
	Rabbit Rabbit `json:"rabbit"`
}

var config ConfigEventHandler

func (conf *ConfigEventHandler) ReadConfig(filename string) int {

	data, err := os.Open(filename)
	if err != nil {
		fmt.Println("Open", err)
		return -1
	}
	defer data.Close()
	decoder := json.NewDecoder(data)
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Decode", err)
		return -1
	}

	fmt.Println(conf)

	return 1
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
}

func connect() (*Consumer, error) {
	var host string
	host = "amqp://" + config.Rabbit.User + ":" + config.Rabbit.Password + "@" + config.Rabbit.Host
	if config.Rabbit.Port != 0 {
		host += ":" + strconv.Itoa(config.Rabbit.Port)
	}
	if config.Rabbit.Vhost != "" {
		host += "/" + config.Rabbit.Vhost
	}
	if config.Rabbit.Heartbeat != 0 {
		host += "?heartbeat=" + strconv.Itoa(config.Rabbit.Heartbeat)
	}

	fmt.Println(host)

	c := &Consumer{
		conn:    nil,
		channel: nil,
		done:    make(chan error),
	}

	var err error
	c.conn, err = amqp.Dial(host)
	if err != nil {
		fmt.Printf("Dial: %s", err)
		return nil, err
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	//defer c.channel.Close()

	fmt.Println("Success connect")
	return c, err
}

func bind(c *Consumer) {
	args := make(amqp.Table)
	args["x-dead-letter-exchange"] = config.Rabbit.Queue.Arguments.Xdeadletterexchange
	args["x-dead-letter-routing-key"] = config.Rabbit.Queue.Arguments.Xdeadletterroutingkey
	args["x-message-ttl"] = config.Rabbit.Queue.Arguments.Xmessagettl
	args["x-expires"] = config.Rabbit.Queue.Arguments.Xexpires

	var err error
	err = c.channel.ExchangeDeclare(
		config.Rabbit.Exchange.Name,
		config.Rabbit.Exchange.Type,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("Exchange Declare: %s", err)
	}

	queue, err := c.channel.QueueDeclare(
		config.Rabbit.Policies+config.Rabbit.Queue.Name,
		false,
		true,
		false,
		false,
		args,
	)
	if err != nil {
		fmt.Printf("Queue Declare: %s", err)
	}

	err = c.channel.QueueBind(
		queue.Name,
		config.Rabbit.Bind.Key,
		config.Rabbit.Exchange.Name,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("Queue Bind: %s", err)
	}
}

func Main() {
	config.ReadConfig("config.json")

	c, err := connect()
	if err != nil {
		fmt.Println(err)
	}

	bind(c)

	msgs, err := c.channel.Consume(
		config.Rabbit.Policies+config.Rabbit.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("Receive Message: %s\n", d.Body)
		}
	}()

	<-forever
}

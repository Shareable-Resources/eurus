package block_cypher_bg

import (
	"encoding/json"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"eurus-backend/service_base/service"
	"eurus-backend/service_base/service_server"

	"github.com/streadway/amqp"
)

type BlockCypherBackground struct {
	service_server.ServiceServer
	completedChannel chan (bool)
	config           *BlockCypherBackgroundConfig
	processor        *BlockCypherBgProcessor
	dbProcessor      *BlockCypherBgDBProcessor
}

func NewBlockCyhpherBackground() *BlockCypherBackground {
	bg := new(BlockCypherBackground)
	bg.config = new(BlockCypherBackgroundConfig)
	bg.ServerConfig = &bg.config.ServerConfigBase
	bg.completedChannel = make(chan bool)

	return bg
}

func (me *BlockCypherBackground) LoadConfig(args *service.ServiceCommandLineArgs) error {
	return me.LoadConfigWithSetting(args, me.config, false, false)
}

func (me *BlockCypherBackground) InitAll() <-chan (bool) {
	me.ServiceServer.InitAuth(me.processInit, me.configEventReceived)
	return me.completedChannel
}

func (me *BlockCypherBackground) processInit(authClient auth_base.IAuth) {
	_, err := me.QueryConfigServer(me.config)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("Unable to get config from config server: ", err)
		panic("Unable to get config from config server: " + err.Error())
	}

	if me.config.CoinListJson != "" {
		err = json.Unmarshal([]byte(me.config.CoinListJson), &me.config.CoinList)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Unmarshal coin json list failed: ", err)
			panic("Unmarshal coin json list failed: " + err.Error())
		}
	}

	err = me.InitDBFromConfig(me.ServerConfig)
	if err != nil {
		log.GetLogger(log.Name.Root).Panicln("Unable to init database: ", err)
	}

	me.dbProcessor = NewDbProcessor(me.DefaultDatabase, &me.config.ServerConfigBase)
	me.processor = NewBlockCypherBgProcessor(me.dbProcessor, me.config)

	me.processor.processUpdateToken()

	me.completedChannel <- true
}

func (me *BlockCypherBackground) configEventReceived(message *amqp.Delivery, topic string, contentType string, content []byte) {

}

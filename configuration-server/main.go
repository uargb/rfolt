package main

import (
	"github.com/BurntSushi/toml"
	"github.com/rocwong/neko"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"rfolt/configuration"
	"strconv"
)

func main() {
	//получим конфиги
	configs := getAllConfigs("./enabled")

	//Запустим http сервер
	server := neko.Classic()

	//Комманда перезагрузки
	server.POST("/jx5xyi5xfr2cv2qxxako1jx8o762r5lt", func(ctx *neko.Context) {
		configs = getAllConfigs("./enabled")
		ctx.Text("OK", http.StatusOK)
	})

	//Роуты конфигов
	server.Group("/:id", func(router *neko.RouterGroup) {
		router.GET("/main", func(ctx *neko.Context) {
			ctx.Json(configs[ctx.Params.ByGet("id")].Main, http.StatusOK)
		})
		router.GET("/network", func(ctx *neko.Context) {
			ctx.Json(configs[ctx.Params.ByGet("id")].Network, http.StatusOK)
		})
		router.GET("/features", func(ctx *neko.Context) {
			ctx.Json(configs[ctx.Params.ByGet("id")].Features, http.StatusOK)
		})
	})

	server.Run("127.0.0.1:9411")
}

func getAllConfigs(dir string) map[string]configuration.Config {
	result := make(map[string]configuration.Config)

	//Просканируем папку конфигов
	configFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	//Распарсим конфиги
	for _, configFile := range configFiles {
		//проверим расширение
		if filepath.Ext(configFile.Name()) != ".toml" {
			continue
		}

		//прочитаем файл конфигурации
		configString, err := ioutil.ReadFile(dir + "/" + configFile.Name())
		if err != nil {
			log.Fatal(err)
		}

		//распарсим его
		var config configuration.Config
		_, err = toml.Decode(string(configString), &config)
		if err != nil {
			log.Fatal(err)
		}

		result[strconv.Itoa(config.Id)] = config
	}

	log.Printf("Загружено конфигураций: %v", len(result))
	return result
}
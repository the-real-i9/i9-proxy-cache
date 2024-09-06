package helpers

import (
	"bufio"
	"i9pxc/src/db"
	"log"
	"os"
	"strings"
)

func loadEnv() error {
	dotenv, err := os.Open(".env")
	if err != nil {
		return err
	}

	env := bufio.NewScanner(dotenv)

	for env.Scan() {
		key, value, found := strings.Cut(env.Text(), "=")
		if !found || strings.HasPrefix(key, "#") {
			continue
		}

		err := os.Setenv(key, value)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func ServerInits() error {
	err := loadEnv()
	if err != nil {
		return err
	}

	db.InitRedisClient()

	return nil
}

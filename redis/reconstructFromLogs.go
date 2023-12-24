package redis

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func (redis *Redis) reconstructFromLogs() error {
	fileScanner := bufio.NewScanner(redis.logFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		commands := strings.Split(fileScanner.Text(), ", ")
		if len(commands) != 4 {
			return errors.New(fmt.Sprint("Failed to parse logs, commands len != 4: ", commands))
		}
		// cmd := commands[0]
		strTimestamp := commands[1]
		strUUID := commands[2]
		strInput := commands[3]

		var data map[string]interface{}
		err := json.Unmarshal([]byte(strInput), &data)
		if err != nil {
			return err
		}

		for key, value := range data {
			parsedUUID, err := uuid.FromString(strUUID)
			if err != nil {
				return err
			}

			timestamp, err := time.Parse(time.RFC3339, strTimestamp)
			redis.Put(key, value, parsedUUID, timestamp)
		}

		fileLines = append(fileLines, fileScanner.Text())
	}

	return nil
}

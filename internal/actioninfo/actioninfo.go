package actioninfo

import (
	"fmt"
	"log"
)

type DataParser interface {
	Parse(datastring string) error
	ActionInfo() (string, error)
}

func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		if err := dp.Parse(data); err != nil {
			log.Printf("failed to parse %q: %v", data, err)
			continue
		}
		info, err := dp.ActionInfo()
		if err != nil {
			log.Printf("failed to format info for %q: %v", data, err)
			continue
		}
		fmt.Println(info)
	}
}

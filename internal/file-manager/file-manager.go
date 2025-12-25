package fm

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Archiker-715/expense-tracker/internal/constants"
)

func CheckExist(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Fatalf("check file exists error: %v", err)
	}

	return true
}

func Create(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("create file error: %w", err)
	}
	return file, nil
}

func Open(fileName string, flag int) (*os.File, error) {
	file, err := os.OpenFile(fileName, flag, 0644)
	if err != nil {
		return nil, fmt.Errorf("open file error: %w", err)
	}
	return file, nil
}

func Write(file *os.File, flag int, input interface{}) error {

	fileName := file.Name()

	if strings.Contains(fileName, "csv") {
		body, ok := input.([][]string)
		if !ok {
			return errors.New("csv writing error, input is not [][]string")
		}

		w := csv.NewWriter(file)
		if flag == os.O_APPEND {
			err := w.Write(body[0])
			if err != nil {
				return fmt.Errorf("writing err: %q", err)
			}
		}

		if flag == os.O_RDWR {
			if err := file.Truncate(0); err != nil {
				return fmt.Errorf("truncate err: %q", err)
			}
			if _, err := file.Seek(0, io.SeekStart); err != nil {
				return fmt.Errorf("seek err: %q", err)
			}
			if err := w.WriteAll(body); err != nil {
				return fmt.Errorf("writing all err: %q", err)
			}
		}

		w.Flush()
		if err := w.Error(); err != nil {
			return fmt.Errorf("error flush data in file: %w", err)
		}
	}

	if fileName == constants.OptionsFileName {
		b, ok := input.([]byte)
		if !ok {
			return errors.New("json writing error, input is not []byte")
		}
		if err := file.Truncate(0); err != nil {
			log.Fatalf("truncate error: %v", err)
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			log.Fatalf("seek error: %v", err)
		}
		if _, err := file.Write(b); err != nil {
			log.Fatalf("writing expense error: %v", err)
		}
	}

	return nil
}

func Read(file *os.File) ([][]string, error) {
	r := csv.NewReader(file)
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek err: %q", err)
	}
	s, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read file error: %q", err)
	}
	if len(s) == 0 {
		return nil, fmt.Errorf("file is empty")
	}
	return s, nil
}

func ReadJson(file *os.File) []byte {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("read file error: %v", err)
	}

	return fileContent
}

func Print(s [][]string) {
	if len(s) == 0 {
		log.Fatalf("file is empty")
	}

	for _, innerS := range s {
		fmt.Println("")
		for i := 0; i < len(innerS); i++ {
			fmt.Printf("%s ", innerS[i])
		}
	}
	fmt.Println("")
}

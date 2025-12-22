package fm

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
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

func Write(file *os.File, flag int, input [][]string) error {
	w := csv.NewWriter(file)

	if flag == os.O_APPEND {
		err := w.Write(input[0])
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
		if err := w.WriteAll(input); err != nil {
			return fmt.Errorf("writing all err: %q", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("error flush data in file: %w", err)
	}

	return nil
}

func Read(file *os.File) ([][]string, error) {
	r := csv.NewReader(file)
	s, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read file error: %q", err)
	}
	if len(s) == 0 {
		return nil, fmt.Errorf("file is empty")
	}
	return s, nil
}

func Print(s [][]string) error {
	if len(s) == 0 {
		return fmt.Errorf("file is empty")
	}

	for _, innerS := range s {
		fmt.Println("")
		for i := 0; i < len(innerS); i++ {
			fmt.Printf("%s ", innerS[i])
		}
	}

	return nil
}

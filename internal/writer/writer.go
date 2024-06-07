package writer

import (
	"os"
	"strconv"
	"sync"
)

type Writer struct {
	Mutex *sync.Mutex
	File  *os.File
}

func MakeWriter() *Writer {
	return &Writer{
		Mutex: new(sync.Mutex),
	}
}

func (writer *Writer) CloseFile() {
	writer.File.Close()
}

func (writer *Writer) CreateCsvIfNotExist(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	writer.File = file
	writer.File.Write([]byte("user_id,minimum_balance,maximum_balance,average_balance,start_timestamp\n"))
	return nil
}

func ConvertFloatToString(value float64) []byte {
	return []byte(strconv.FormatFloat(value, 'f', 4, 64))
}

func ConvertIntToString(value int) []byte {
	return []byte(strconv.Itoa(value))
}

func (writer Writer) Write(a string, b float64, c float64, d float64, e int) {
	writer.Mutex.Lock()
	writer.File.Write([]byte(a))
	writer.File.Write([]byte{','})
	writer.File.Write(ConvertFloatToString(b))
	writer.File.Write([]byte{','})
	writer.File.Write(ConvertFloatToString(c))
	writer.File.Write([]byte{','})
	writer.File.Write(ConvertFloatToString(d))
	writer.File.Write([]byte{','})
	writer.File.Write(ConvertIntToString(e))
	writer.File.Write([]byte("\n"))
	writer.Mutex.Unlock()
}

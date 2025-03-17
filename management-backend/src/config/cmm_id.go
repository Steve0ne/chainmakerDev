package config

import (
	"bufio"
	"io"
	"os"

	"chainmaker.org/chainmaker/common/v2/random/uuid"
)

var CMMID, _ = GetCMMId()

const IdFile = "CMMID"

func GetCMMId() (string, error) {

	f, err := os.OpenFile(IdFile, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	buf := make([]byte, 32)

	r := bufio.NewReader(f)
	n, err := r.Read(buf)

	if err != nil && err != io.EOF {
		panic(err)
	}

	id := string(buf[:n])

	if n < 5 {
		id = uuid.GetUUID()
		_, err = f.WriteAt([]byte(id), 0)
	}

	return id, nil
}

package uuid

import (
	"crypto/rand"
	"encoding/hex"
)

func New() string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// uuid версии 4 состоит из 122 случайных бит и 6 бит, обозначающих версию и варианты
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx-0100xxxxxxxxxxxx-10xxxxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// биты версии (соответствует седьмому байту в наборе случайных байт) стираются и устанавливаются в 0100 (бинарное 4)
	randomBytes[6] = (randomBytes[6] & 0b00001111) | 0b01000000
	// биты варианта (соответствует девятому байту) стираются и устанавливаются в 10 (hex-репрезентация битов 10xx будет принимать значения 8..d)
	randomBytes[8] = (randomBytes[8] & 0b00111111) | 0b10000000

	uuid := make([]byte, 36)
	hex.Encode(uuid[0:8], randomBytes[0:4])
	uuid[8] = '-'
	hex.Encode(uuid[9:13], randomBytes[4:6])
	uuid[13] = '-'
	hex.Encode(uuid[14:18], randomBytes[6:8])
	uuid[18] = '-'
	hex.Encode(uuid[19:23], randomBytes[8:10])
	uuid[23] = '-'
	hex.Encode(uuid[24:], randomBytes[10:16])

	return string(uuid)
}

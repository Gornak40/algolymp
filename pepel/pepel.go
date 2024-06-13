package pepel

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type MagicPair struct {
	InfSHA256 string
	AnsBase64 string
}

func GenMagicPair(inf, ans string, zFlag bool) (*MagicPair, error) {
	finf, err := os.Open(inf)
	if err != nil {
		return nil, err
	}
	defer finf.Close()

	hshSum := sha256.New()
	if _, err := io.Copy(hshSum, finf); err != nil {
		logrus.WithError(err).Fatal("failed to calculate hash")
	}
	hsh := hex.EncodeToString(hshSum.Sum(nil))

	fans, err := os.Open(ans)
	if err != nil {
		logrus.WithError(err).Fatal("failed to open answer file")
	}
	defer fans.Close()
	ansData, err := io.ReadAll(fans)
	if err != nil {
		logrus.WithError(err).Fatal("failed to read answer file")
	}

	if zFlag {
		var b bytes.Buffer
		w := zlib.NewWriter(&b)
		if _, err := w.Write(ansData); err != nil {
			return nil, err
		}
		w.Close()
		ansData = b.Bytes()
	}

	ansEnc := base64.StdEncoding.EncodeToString(ansData)

	return &MagicPair{
		InfSHA256: hsh,
		AnsBase64: ansEnc,
	}, nil
}

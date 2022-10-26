package idproto

import (
	"crypto/sha256"
	"encoding/base64"
	"path/filepath"
)

func MailTopicDirpath(topic string) string {
	return filepath.Join(IdentityRoot, "mail", topicHash(topic))
}

func topicHash(topic string) string {
	h := sha256.New()
	if _, err := h.Write([]byte(topic)); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

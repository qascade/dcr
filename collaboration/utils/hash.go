package utils

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strconv"
)

func HashBool(b bool) string {
	return HashString(strconv.FormatBool(b))
}

func HashInt(i int) string {
	return HashString(strconv.Itoa(i))
}

func HashString(str string) string {
	return HashBytes([]byte(str))
}

func HashBytes(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:])[0:8]
}

func HashStringMap(stringMap map[string]string) string {
	h := sha256.New()
	keys := make([]string, len(stringMap))
	i := 0
	for k := range stringMap {
		keys[i] = k
		i += 1
	}
	sort.Strings(keys)
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte(stringMap[k]))
	}
	sum := h.Sum(nil)
	hashStr := fmt.Sprintf("%x", sum)[0:8]
	return hashStr
}

func HashUnorderedStringList(stringList []string) string {
	sort.Strings(stringList)
	return HashOrderedStringList(stringList)
}

func HashOrderedStringList(stringList []string) string {
	h := sha256.New()
	for _, s := range stringList {
		h.Write([]byte(s))
	}
	sum := h.Sum(nil)
	hashStr := fmt.Sprintf("%x", sum)[0:8]
	return hashStr
}

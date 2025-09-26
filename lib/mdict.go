package lib

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("default")

type Mdict struct {
	*MdictBase
}

func New(filename string) (*Mdict, error) {
	dictType := MdictTypeMdx
	if strings.ToLower(filepath.Ext(filename)) == ".mdd" {
		dictType = MdictTypeMdd
	}

	mdict := &Mdict{
		MdictBase: &MdictBase{
			FilePath: filename,
			FileType: dictType,
		},
	}
	return mdict, mdict.init()
}

func (mdict *Mdict) init() error {
	// 读取词典头
	err := mdict.readDictHeader()
	if err != nil {
		return err
	}

	// 读取 key block 元信息
	err = mdict.readKeyBlockMeta()
	if err != nil {
		return err
	}

	return nil
}

// LoadMetadata loads all metadata except the key entries themselves.
// This is a fast operation used to prepare a parser for lookups when the index is already built.
func (mdict *Mdict) LoadMetadata() error {
	err := mdict.readKeyBlockInfo()
	if err != nil {
		return err
	}

	err = mdict.readRecordBlockMeta()
	if err != nil {
		return err
	}

	err = mdict.readRecordBlockInfo()
	if err != nil {
		return err
	}

	return nil
}

// BuildIndex 构建索引
func (mdict *Mdict) BuildIndex() error {
	err := mdict.readKeyBlockInfo()
	if err != nil {
		return err
	}

	err = mdict.readKeyEntries()
	if err != nil {
		return err
	}

	err = mdict.readRecordBlockMeta()
	if err != nil {
		return err
	}

	err = mdict.readRecordBlockInfo()
	if err != nil {
		return err
	}

	return nil
}

func (mdict *Mdict) LocateByKeywordEntry(entry *MDictKeywordEntry) ([]byte, error) {
	if entry == nil {
		return nil, errors.New("invalid mdict keyword entry")
	}
	return mdict.MdictBase.LocateByKeywordEntry(entry)
}

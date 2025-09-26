package lib

type MdictType int

const (
	MdictTypeMdd MdictType = 1
	MdictTypeMdx MdictType = 2

	EncryptNoEnc      = 0
	EncryptRecordEnc  = 1
	EncryptKeyInfoEnc = 2
	NumfmtBe8bytesq   = 0
	NumfmtBe4bytesi   = 1
	EncodingUtf8      = 0
	EncodingUtf16     = 1
	EncodingBig5      = 2
	ENCODING_GBK      = 3
	ENCODING_GB2312   = 4
	EncodingGb18030   = 5
)

type MdictBase struct {
	FilePath string
	FileType MdictType
	Meta     *MdictMeta

	header       *mdictHeader
	keyBlockMeta *mdictKeyBlockMeta
	keyBlockInfo *mdictKeyBlockInfo
	KeyBlockData *MdictKeyBlockData

	recordBlockMeta *mdictRecordBlockMeta
	RecordBlockInfo *MdictRecordBlockInfo
	//RecordBlockData *MDictRecordBlockData
}

/********************************
 *    private data type          *
 ********************************/
type mdictHeader struct {
	headerBytesSize          uint32
	headerInfoBytes          []byte
	headerInfo               string
	adler32Checksum          uint32
	dictionaryHeaderByteSize int64
}

type MdictMeta struct {
	EncryptType  int
	Version      float32
	NumberWidth  int
	NumberFormat int
	Encoding     int

	// key-block part bytes start offset in the mdx/mdd file
	KeyBlockMetaStartOffset int64

	Description              string
	Title                    string
	CreationDate             string
	GeneratedByEngineVersion string
}

type mdictKeyBlockMeta struct {
	// keyBlockNum key block number size
	keyBlockNum int64
	// entriesNums entries number size
	entriesNum int64
	// key-block information size (decompressed)
	keyBlockInfoDecompressSize int64
	// key-block information size (compressed)
	keyBlockInfoCompressedSize int64
	// key-block Data Size (decompressed)
	keyBlockDataTotalSize int64
	// key-block information start position in the mdx/mdd file
	keyBlockInfoStartOffset int64
}

type mdictKeyBlockInfo struct {
	keyBlockEntriesStartOffset int64
	keyBlockInfoList           []*mdictKeyBlockInfoItem
}

type mdictKeyBlockInfoItem struct {
	firstKey                      string
	firstKeySize                  int
	lastKey                       string
	lastKeySize                   int
	keyBlockInfoIndex             int
	keyBlockCompressSize          int64
	keyBlockCompAccumulator       int64
	keyBlockDeCompressSize        int64
	keyBlockDeCompressAccumulator int64
}

type MdictKeyBlockData struct {
	KeyEntries                 []*MDictKeywordEntry
	KeyEntriesSize             int64
	RecordBlockMetaStartOffset int64
}

type mdictRecordBlockMeta struct {
	keyRecordMetaStartOffset int64
	keyRecordMetaEndOffset   int64

	recordBlockNum          int64
	entriesNum              int64
	recordBlockInfoCompSize int64
	recordBlockCompSize     int64
}
type MdictRecordBlockInfo struct {
	RecordInfoList             []*MdictRecordBlockInfoListItem
	RecordBlockInfoStartOffset int64
	RecordBlockInfoEndOffset   int64
	RecordBlockDataStartOffset int64
}

type MdictRecordBlockInfoListItem struct {
	CompressSize                int64
	DeCompressSize              int64
	CompressAccumulatorOffset   int64
	DeCompressAccumulatorOffset int64
}

/********************************
 *    public data type          *
 ********************************/

type MDictKeywordEntry struct {
	RecordStartOffset int64
	RecordEndOffset   int64
	KeyWord           string
	KeyBlockIdx       int64
}

type MDictKeywordIndex struct {
	//encoding                            int
	//encryptType                         int
	KeywordEntry MDictKeywordEntry
	RecordBlock  MDictKeywordIndexRecordBlock
}

type MDictKeywordIndexRecordBlock struct {
	DataStartOffset          int64
	CompressSize             int64
	DeCompressSize           int64
	KeyWordPartStartOffset   int64
	KeyWordPartDataEndOffset int64
}

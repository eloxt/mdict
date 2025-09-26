package replacer

type Replacer interface {
	Replace(dictId string, htmlContent string) string
}

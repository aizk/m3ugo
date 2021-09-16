package fetch

// Fetcher m3u8 file from url
type Fetcher interface {
	Fetch(url string) (mu3Url, fileName string, err error)
}

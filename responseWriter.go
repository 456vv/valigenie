package valigenie
	
import (
	"io"
)
type ResponseWriter interface{
	Ask(reply string, askedInfos []string)
	Askf(reply string, askedInfos []string, args... interface{})
	Confirm(reply string)
	Confirmf(reply string, args... interface{})
	Result(reply string)
	Resultf(reply string, args... interface{})
	Error(errStr, code string)
	WriteTo(w io.Writer) (n int64, err error) 
}
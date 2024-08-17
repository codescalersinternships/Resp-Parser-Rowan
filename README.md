# Resp-Parser-Rowan
A Go library for parsing Redis Serialization Protocol (RESP) data. 

## Resources/Credits:
1. Taking inspo from [this resp library](https://github.com/tidwall/resp?tab=readme-ov-file#server)
2. Documentation in [redis docs](https://redis.io/docs/latest/develop/reference/protocol-spec/)

## Examples:
```
import (
	"bytes"
	"fmt"

	"github.com/codescalersinternships/Resp-Parser-Rowan/resp"
)

func main() {
	var value resp.Value
	
	raw := "*3\r\n:1\r\n:2\r\n:3\r\n"
	rd := resp.NewReader(bytes.NewBufferString(raw))
	value, err := rd.ReadValue()
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range value.Array {
		fmt.Println(value.Integer)
	}
	fmt.Println(value.Type() == resp.Array)
}
```
### Results:
```
1
2
3
true
```
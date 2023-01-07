package smCrypto

import (
	"fmt"
	sm2 "github.com/sunshibao/go-sm/sm2"
	sm4 "github.com/sunshibao/go-sm/sm4"
	"os"
)

var Sm2Crypto *sm2.Sm2Cypher
var Sm4Crypto *sm4.Sm4Cypher
var err error

//
func init() {
	//pub1: 04a984bf2b627d50995e3a167b347d61af6ce1b66dae959fac52f6c6574f9c731c3df2f3ae17b8423b81b06f793850b4f78198f1e50a23b85ac77f5e05a9e672d2
	//prv1: 161e5bca7c543359563f01b1e4158224c39714478a61391bf3c01162eead3290
	//pub2: 044359ae2a8688df666f6d54814706ff033f36a6c6d273d6cb17d2531550f3d48233d8f66e589d7dcd70a6b809d4e17e7c97d2ac24e41eee4ab6fe2c7bd0d485bc
	//prv2: 84bfcd3a08692306dfb4a6f18415ce944023f0c83877750ee6f8e687045a43a5
	Sm2Crypto, err = sm2.NewSm2(sm2.Option{
		Mode:    0,
		Random:  "",
		PubStr:  "04a984bf2b627d50995e3a167b347d61af6ce1b66dae959fac52f6c6574f9c731c3df2f3ae17b8423b81b06f793850b4f78198f1e50a23b85ac77f5e05a9e672d2",
		PrvStr:  "84bfcd3a08692306dfb4a6f18415ce944023f0c83877750ee6f8e687045a43a5",
		KeyType: "hex",
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	Sm4Crypto, err = sm4.NewSm4("@PHXGV9Tb9V+8-J&", "sagws8poOl%s4M3e")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

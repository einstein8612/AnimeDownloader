package decoders

import (
	"io/ioutil"

	"github.com/robertkrimen/otto"
)

var TWIST_SECRET_KEY = "LXgIVP&PorO68Rq7dTx8N^lP!Fa5sGJ^*XK"
var vm *otto.Otto

func twistInit() {
	gibberishBytes, _ := ioutil.ReadFile("./gibberishaes.js")
	vm = otto.New()
	vm.Run(gibberishBytes)
}

func DecodeTwist(source string) string {
	val, _ := vm.Run(`GibberishAES.dec("` + source + `", "` + TWIST_SECRET_KEY + `")`)
	decodedSource, _ := val.ToString()
	return decodedSource
}

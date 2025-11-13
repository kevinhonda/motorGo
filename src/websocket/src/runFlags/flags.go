package runFlags

import "flag"

type FlagsType struct {
	Test    *bool
	Dev     *bool
	Version *bool
}

var Flags *FlagsType

func Parse() {
	test := flag.Bool("t", false, "Testar a conexão")
	dev := flag.Bool("d", false, "Definir como ambiente dev")
	v := flag.Bool("v", false, "Exibe a versão do build no terminal")
	flag.Parse()

	Flags = &FlagsType{
		Dev:     dev,
		Version: v,
		Test:    test,
	}
}

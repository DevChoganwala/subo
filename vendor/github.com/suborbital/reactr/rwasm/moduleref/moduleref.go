package moduleref

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

// WasmModuleRef is a reference to a Wasm module (either its filepath or its bytes)
type WasmModuleRef struct {
	Filepath string
	Name     string
	data     []byte
}

// RefWithData returns a module ref from module bytes
func RefWithData(name string, data []byte) *WasmModuleRef {
	ref := &WasmModuleRef{
		Name: name,
		data: data,
	}

	return ref
}

// Bytes returns the bytes for the module
func (w *WasmModuleRef) Bytes() ([]byte, error) {
	if w.data == nil {
		if w.Filepath == "" {
			return nil, errors.New("missing Wasm module filepath in ref")
		}

		bytes, err := ioutil.ReadFile(w.Filepath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to ReadFile for Wasm module")
		}

		w.data = bytes
	}

	return w.data, nil
}

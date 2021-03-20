package rwasm

import (
	"github.com/pkg/errors"
	"github.com/wasmerio/wasmer-go/wasmer"
)

func getFFIResult() *HostFn {
	fn := func(args ...wasmer.Value) (interface{}, error) {
		pointer := args[0].I32()
		ident := args[1].I32()

		ret := get_ffi_result(pointer, ident)

		return ret, nil
	}

	return newHostFn("get_ffi_result", 2, true, fn)
}

func get_ffi_result(pointer int32, identifier int32) int32 {
	envLock.RLock()
	defer envLock.RUnlock()

	inst, err := instanceForIdentifier(identifier, false)
	if err != nil {
		logger.Error(errors.Wrap(err, "[rwasm] failed to instanceForIdentifier"))
		return -1
	}

	result, err := inst.useFFIResult()
	if err != nil {
		logger.Error(errors.Wrap(err, "[rwasm] failed to useFFIResult"))
		return -1
	}

	inst.writeMemoryAtLocation(pointer, result)

	return 0
}

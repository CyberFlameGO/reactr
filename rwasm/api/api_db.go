package api

import (
	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm/runtime"
)

func DBExecHandler() runtime.HostFn {
	fn := func(args ...interface{}) (interface{}, error) {
		queryType := args[0].(int32)
		namePointer := args[1].(int32)
		nameSize := args[2].(int32)
		ident := args[3].(int32)

		ret := db_exec(queryType, namePointer, nameSize, ident)

		return ret, nil
	}

	return runtime.NewHostFn("db_exec", 4, true, fn)
}

func db_exec(queryType, namePointer, nameSize, identifier int32) int32 {
	inst, err := runtime.InstanceForIdentifier(identifier, false)
	if err != nil {
		runtime.InternalLogger().Error(errors.Wrap(err, "[rwasm] alert: failed to InstanceForIdentifier"))
		return -1
	}

	nameBytes := inst.ReadMemory(namePointer, nameSize)
	name := string(nameBytes)

	vars, err := inst.Ctx().UseVars()
	if err != nil {
		runtime.InternalLogger().Error(errors.Wrap(err, "[rwasm] failed to UseVars"))
	}

	queryResult, err := inst.Ctx().Database.ExecQuery(queryType, name, varsToInterface(vars))
	if err != nil {
		runtime.InternalLogger().ErrorString("[rwasm] failed to ExecQuery", name, err.Error())

		res, _ := inst.Ctx().SetFFIResult(nil, err)
		return res.FFISize()
	}

	res, _ := inst.Ctx().SetFFIResult(queryResult, nil)

	return res.FFISize()
}

func varsToInterface(vars []rt.FFIVariable) []interface{} {
	iVars := []interface{}{}

	for _, v := range vars {
		iVars = append(iVars, v.Value)
	}

	return iVars
}

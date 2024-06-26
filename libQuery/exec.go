package libQuery

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/hmmftg/requestCore/libError"
	"github.com/hmmftg/requestCore/response"
	"github.com/hmmftg/requestCore/webFramework"
)

func (m QueryRunnerModel) SetVariableCommand() string {
	return m.SetVariable
}

const (
	ERROR_CALLING_DB_FUNCTION = "ERROR_CALLING_DB_FUNCTION"
)

func (m QueryRunnerModel) CallDbFunction(callString string, args ...any) (int, string, error) {
	errPing := m.DB.Ping()
	if errPing != nil {
		log.Println("error in ping", errPing)
	}
	_, err := m.DB.Exec(callString, args...)
	if err != nil {
		return -3, ERROR_CALLING_DB_FUNCTION, libError.Join(err, "CallDbFunction[Exec](%s,%v)", callString, args)
	}

	return 0, "OK", nil
}

const (
	QueryCheckNotExists DmlCommandType = iota
	QueryCheckExists
	Insert
	Update
	Delete
)

func (command DmlCommand) Execute(core QueryRunnerInterface, moduleName, methodName string) (any, response.ErrorState) {
	return command.ExecuteWithContext(nil, context.Background(), moduleName, methodName, core)
}

func GetDmlResult(resultDb sql.Result, rows map[string]string) DmlResult {
	resp := DmlResult{
		Rows: rows,
	}
	resp.LastInsertId, _ = resultDb.LastInsertId()
	resp.RowsAffected, _ = resultDb.RowsAffected()
	return resp
}

func GetLocalArgs(parser webFramework.RequestParser, args []any) []any {
	result := make([]any, len(args))
	for id, arg := range args {
		namedArg, ok := arg.(sql.NamedArg)
		if ok {
			stringArg, ok := namedArg.Value.(string)
			if ok && strings.HasPrefix(stringArg, "w.local:") {
				parts := strings.Split(stringArg, ":")
				namedArg.Value = parser.GetLocal(parts[1])
			}
			result[id] = namedArg
		} else {
			result[id] = arg
		}
	}
	return result
}

func GetOutArgs(parser webFramework.RequestParser, args ...any) map[string]string {
	rows := map[string]string{}
	for id, arg := range args {
		switch dbParameter := arg.(type) {
		case sql.NamedArg:
			switch namedParameter := dbParameter.Value.(type) {
			case sql.Out:
				if namedParameter.Dest != nil {
					switch outValue := namedParameter.Dest.(type) {
					case string:
						rows[dbParameter.Name] = outValue
					case *string:
						rows[dbParameter.Name] = *outValue
					case int64:
						rows[dbParameter.Name] = fmt.Sprintf("%d", outValue)
					case *int64:
						rows[dbParameter.Name] = fmt.Sprintf("%d", *outValue)
					default:
						log.Printf("wrong db-out parameter type %T\n", namedParameter.Dest)
					}
					parser.SetLocal(dbParameter.Name, rows[dbParameter.Name])
				}
			}
		case sql.Out:
			if dbParameter.Dest != nil {
				name := fmt.Sprintf("not named arg %d", id)
				switch outValue := dbParameter.Dest.(type) {
				case string:
					rows[name] = outValue
				case *string:
					rows[name] = *outValue
				case int64:
					rows[name] = fmt.Sprintf("%d", outValue)
				case *int64:
					rows[name] = fmt.Sprintf("%d", *outValue)
				default:
					log.Printf("wrong db-out parameter type %T\n", dbParameter.Dest)
				}
				parser.SetLocal(name, rows[name])
			}
		}
	}
	return rows
}

func (command DmlCommand) ExecuteWithContext(parser webFramework.RequestParser, ctx context.Context, moduleName, methodName string, core QueryRunnerInterface) (any, response.ErrorState) {
	switch command.Type {
	case QueryCheckExists:
		_, desc, data, resp, err := CallSql[QueryData](command.Command, core, GetLocalArgs(parser, command.Args)...)
		if err != nil {
			if command.CustomError != nil {
				return nil, command.CustomError
			}
			return nil, response.ToError(desc, data, libError.Join(err, "CheckExists: %s", command.Name))
		}
		if desc == NO_DATA_FOUND {
			return nil, response.ToError(NO_DATA_FOUND, NO_DATA_FOUND_DESC, fmt.Errorf("CheckExists: %s=> %s", command.Name, NO_DATA_FOUND))
		}
		return resp, nil
	case QueryCheckNotExists:
		_, desc, data, resp, err := CallSql[QueryData](command.Command, core, GetLocalArgs(parser, command.Args)...)
		if len(desc) > 0 && desc == NO_DATA_FOUND && resp == nil {
			return nil, nil // OK
		}
		if err != nil {
			if command.CustomError != nil {
				return nil, command.CustomError
			}
			return nil, response.ToError(desc, data, libError.Join(err, "CheckNotExists: %s", command.Name))
		}
		return nil, response.ToError(DUPLICATE_FOUND, DUPLICATE_FOUND_DESC, fmt.Errorf("CheckNotExists: %s=> %s", command.Name, DUPLICATE_FOUND))
	case Insert:
		resp, err := core.Dml(ctx, moduleName, methodName, command.Command, GetLocalArgs(parser, command.Args)...)
		if err != nil {
			if command.CustomError != nil {
				return nil, command.CustomError
			}
			return nil, response.ToError(ErrorExecuteDML, nil, libError.Join(err, "%s: %s", command.Type, command.Name))
		}
		outValues := GetOutArgs(parser, command.Args...)
		return GetDmlResult(resp, outValues), nil
	case Update:
		resp, err := core.Dml(ctx, moduleName, methodName, command.Command, GetLocalArgs(parser, command.Args)...)
		if err != nil {
			if command.CustomError != nil {
				return nil, command.CustomError
			}
			return nil, response.ToError(ErrorExecuteDML, nil, libError.Join(err, "%s: %s", command.Type, command.Name))
		}
		outValues := GetOutArgs(parser, command.Args...)
		return GetDmlResult(resp, outValues), nil
	case Delete:
		resp, err := core.Dml(ctx, moduleName, methodName, command.Command, GetLocalArgs(parser, command.Args)...)
		if err != nil {
			if command.CustomError != nil {
				return nil, command.CustomError
			}
			return nil, response.ToError(ErrorExecuteDML, nil, libError.Join(err, "%s: %s", command.Type, command.Name))
		}
		outValues := GetOutArgs(parser, command.Args...)
		return GetDmlResult(resp, outValues), nil
	}
	return nil, nil
}

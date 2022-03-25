package service

import (
	"context"
	"errors"
	"readingo/constant"
	"readingo/dao"
	"readingo/model"
	"sort"
	"strings"
)

func OperateRedis(ctx context.Context, req interface{}) (resp interface{}, code int, err error) {
	request := req.(*model.OperateRedisReq)

	request.Action = strings.ToUpper(request.Action)
	role := GetCurrentUserRole(ctx)
	if role == "" {
		code = model.NoPermission
		return
	}
	if argsFormat := getArgsFormatByRole(role, request.Action); argsFormat != "" {
		args := populateArgs(request.Action, argsFormat, request)
		if len(args) <= 0 {
			return nil, model.WrongArgs, errors.New("wrong format in server")
		}

		if request.Type == "redisc" {
			resp, err = executeCommandInRedisc(ctx, request.Host, args)
		} else {
			resp, err = executeCommandInRedis(ctx, request.Host, request.Index, args)
		}

		if err != nil {
			code = model.Error
		}
	} else {
		return nil, model.UnsupportedAction, errors.New("unsupported action:" + request.Action)
	}
	return
}

func RefreshDBTree(ctx context.Context, req interface{}) (resp interface{}, code int, err error) {
	if err = dao.RefreshDBTree(ctx); err != nil {
		code = model.Error
	} else {
		resp = model.GetDBTreeResp{Hosts: getAllDBHosts(ctx), LastUpdateTime: dao.GetLastUpdateTime()}
	}
	return
}

func GetDBTree(ctx context.Context, req interface{}) (resp interface{}, code int, err error) {
	return model.GetDBTreeResp{Hosts: getAllDBHosts(ctx), LastUpdateTime: dao.GetLastUpdateTime()}, 0, nil
}

func GetSupportedActions(ctx context.Context, req interface{}) (resp interface{}, code int, err error) {
	return model.GetSupportedActionResp{Actions: getCommandsByRole(GetCurrentUserRole(ctx))}, 0, nil
}

func getAllDBHosts(ctx context.Context) []model.HostInfo {
	hosts := make([]model.HostInfo, 0)
	for _, node := range dao.GetHostsCache(ctx) {
		hosts = append(hosts, node)
	}
	for name, _ := range dao.GetAllRedisClusterClients() {
		hosts = append(hosts, model.HostInfo{
			Host:       name,
			Type:       "redisc",
			Partitions: []model.PartitionInfo{{Index: 0, DBSize: 0}},
		})
	}
	sort.Slice(hosts, func(i, j int) bool {
		return strings.Compare(strings.ToUpper(hosts[i].Host), strings.ToUpper(hosts[j].Host)) < 0
	})
	return hosts
}

func executeCommandInRedis(ctx context.Context, host, idx string, args []interface{}) (interface{}, error) {
	conn, err := dao.GetConn(host, idx)
	if err != nil {
		return nil, err
	}
	reply, err := conn.Do(ctx, args...).Result()
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func executeCommandInRedisc(ctx context.Context, host string, args []interface{}) (interface{}, error) {
	conn, err := dao.GetClusterConn(host)
	if err != nil {
		return nil, err
	}
	reply, err := conn.Do(ctx, args...).Result()
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// populateArgs TODO 之后用反射来做
func populateArgs(action string, argsFormat string, req *model.OperateRedisReq) []interface{} {
	args := []interface{}{action}

	ps := splitCommandFormat(argsFormat)
	for _, p := range ps {
		switch p {
		case "*Key":
			args = append(args, req.Key)
		case "*Field":
			args = append(args, req.Field)
		case "*Value":
			args = append(args, req.Value)
		case "*Score":
			args = append(args, req.Score)
		case "*Min":
			args = append(args, req.Min)
		case "*Max":
			args = append(args, req.Max)
		case "[EX *TTL]":
			if req.TTL != "" {
				args = append(args, "EX "+req.TTL)
			}
		default:
			args = make([]interface{}, 0)
			break
		}
	}
	return args
}

func splitCommandFormat(format string) []string {
	rs := []rune(format)
	rs = append(rs, ' ')
	ps := make([]string, 0)
	startIndex := 0
	inClosure := false
	for i := 0; i < len(rs); i++ {
		if rs[i] == ' ' && !inClosure {
			ps = append(ps, string(rs[startIndex:i]))
			startIndex = i + 1
		} else if rs[i] == '[' {
			inClosure = true
		} else if rs[i] == ']' {
			inClosure = false
		}
	}
	// 有中括号没有闭合时，视为格式定义有误
	if inClosure {
		return []string{}
	}
	return ps
}

func extractRequiredParams(raw string) []string {
	params := make([]string, 0)

	for _, p := range strings.Split(raw, " ") {
		switch p {
		case "*Key":
			params = append(params, "Key")
		case "*Field":
			params = append(params, "Field")
		case "*Value":
			params = append(params, "Value")
		case "*Min":
			params = append(params, "Min")
		case "*Max":
			params = append(params, "Max")
		case "*Score":
			params = append(params, "Score")
		case "[EX *TTL]":
			params = append(params, "TTL")
		}
	}
	return params
}

func getArgsFormatByRole(role, action string) string {
	if argsFormat, ok := constant.SupportedReadCommands[action]; ok {
		return argsFormat
	}

	if role == constant.RoleReadWrite {
		if argsFormat, ok := constant.SupportedWriteCommands[action]; ok {
			return argsFormat
		}
	}
	return ""
}

func getCommandsByRole(role string) []model.Action {
	actions := make([]model.Action, 0)
	for k, v := range constant.SupportedReadCommands {
		actions = append(actions, model.Action{
			Action:         k,
			RequiredParams: extractRequiredParams(v),
			Tips:           v,
		})
	}

	if role == constant.RoleReadWrite {
		for k, v := range constant.SupportedWriteCommands {
			actions = append(actions, model.Action{
				Action:         k,
				RequiredParams: extractRequiredParams(v),
				Tips:           v,
			})
		}
	}

	sort.Slice(actions, func(i, j int) bool {
		return strings.Compare(strings.ToUpper(actions[i].Action), strings.ToUpper(actions[j].Action)) < 0
	})
	return actions
}

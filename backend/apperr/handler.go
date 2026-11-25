package apperr

import (
	"fmt"
	"net/http"

	"deepseek-backend/response"
)

/*
	错误处理逻辑
	失败：业务代码中返回错误→调用本包handleError函数→调用apperr包HandleError函数→调用response包的Fail
	成功：业务代码→调用response包的Success
*/

func HandleError(w http.ResponseWriter, r *http.Request, err error, business HTTPErrorMap) {
	// 看错误是否输入参数的business，属于就返回映射体
	he, ok := lookup(err, httpErrorMap)
	if ok {
		write(w, r, he)
		return
	}

	// 看错误是否属于本包的httpErrorMap，属于就返回映射体
	he, ok = lookup(err, httpErrorMap)
	if ok {
		write(w, r, he)
		return
	}

	// 兜底：按照本包中的ErrInternal错误处理
	write(w, r, httpErrorMap[ErrInternal])

	// 打印，引起开发者注意
	fmt.Printf("unknown error: %v\n", err)
}

func lookup(err error, heMap HTTPErrorMap) (HTTPError, bool) {
	he, ok := heMap[err]
	if ok {
		return he, true
	}

	return HTTPError{}, false
}

// he是http error的缩写
func write(w http.ResponseWriter, r *http.Request, he HTTPError) {
	response.Fail(w, r, he.Status, string(he.Code), he.Message)
}

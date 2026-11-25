package store

/*
	store包是session包的附属包
	理想状态下，这些错误应该定义在session包
	但为了精细处理（更早返回自定义错误）需要在store包就用到这些错误
	为了避免循环引入就定义在这里了
*/

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrContentNotFound = errors.New("content block not found")
)

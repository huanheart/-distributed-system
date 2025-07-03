package controller

import "MyChat/common/code"

type MusicDetail struct {
	FilePath  string `json:"file_path" binding:"required"`
	LikeCount int64  `json:"like_count"binding:"required"`
}

type Response struct {
	StatusCode code.Code `json:"status_code"`
	StatusMsg  string    `json:"status_msg,omitempty"`
}

func (r *Response) CodeOf(code code.Code) Response {
	if nil == r {
		r = new(Response)
	}
	r.StatusCode = code
	r.StatusMsg = code.Msg()
	return *r
}

func (r *Response) Success() {
	r.CodeOf(code.CodeSuccess)
}

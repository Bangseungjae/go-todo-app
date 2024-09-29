package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
*
omitempty는 Go 언어의 struct 태그 중 하나로,
JSON 직렬화 과정에서 특정 필드가 비어있거나 기본값인 경우 해당 필드를 JSON 출력에서 생략하도록 지시합니다.
이 태그는 특히 API 응답을 구성할 때 유용합니다.
*/
type ErrResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("encode response error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		rsp := ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("write error response error: %v", err)
		}
		return
	}
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}

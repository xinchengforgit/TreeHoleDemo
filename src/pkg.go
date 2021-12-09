package main

type UniverseResponse struct {
	Success bool        `json: "successs"`
	Hint    string      `json: "hint"`
	Data    interface{} `json: "data"`
}

//universe return
func ErrorResponse(hint string) UniverseResponse {
	return UniverseResponse{
		Success: false,
		Hint:    hint,
		Data:    nil,
	}
}

func SuccessResponse(data interface{}) UniverseResponse {
	return UniverseResponse{
		Success: true,
		Hint:    "",
		Data:    data,
	}
}

package models

type CalculateRequest struct {
    Expression string `json:"expression" binding:"required"`
}

type TaskResultRequest struct {
    ID     int     `json:"id"`
    Result float64 `json:"result"`
}

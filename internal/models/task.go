package models

type Task struct {
    ID            int     `json:"id"`
    ExpressionID  string  `json:"expression_id"`
    Arg1          string  `json:"arg1"`
    Arg2          string  `json:"arg2,omitempty"`
    Operation     string  `json:"operation"`
    OperationTime int     `json:"operation_time"`
    Status        string  `json:"status"`
}

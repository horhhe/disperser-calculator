package config

import (
    "os"
    "strconv"
)

type EnvConfig struct {
    TimeAddition       int
    TimeSubtraction    int
    TimeMultiplication int
    TimeDivision       int
    TimeEvaluation     int
}

func InitEnv() *EnvConfig {
    return &EnvConfig{
        TimeAddition:       getIntFromEnv("TIME_ADDITION_MS", 1000),
        TimeSubtraction:    getIntFromEnv("TIME_SUBTRACTION_MS", 1000),
        TimeMultiplication: getIntFromEnv("TIME_MULTIPLICATION_MS", 1000),
        TimeDivision:       getIntFromEnv("TIME_DIVISION_MS", 1000),
        TimeEvaluation: getIntFromEnv("TIME_EVALUATION_MS", 2000),
    }
}

func getIntFromEnv(envKey string, defaultVal int) int {
    if val := os.Getenv(envKey); val != "" {
        if i, err := strconv.Atoi(val); err == nil {
            return i
        }
    }
    return defaultVal
}

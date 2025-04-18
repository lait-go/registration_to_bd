package LogWork

import (
	conf "back/config"
	Error "back/internal/err"
	"back/internal/utils"
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

type LogStruct struct {
	Logger slog.Logger
}

type readLog struct {
	Log   string  `json:"log"`
	Time  string  `json:"time"`
	Level string  `json:"level"`
}

func (logStruct *LogStruct) Info(message string) {
	logStruct.Logger.Info(message)

	pars := readLog{
		Log: message,
		Time: utils.TimeFormat(),
		Level: "INFO",
	}

	ReadLog(pars)
}

func (logStruct *LogStruct) Debug(message string) {
	logStruct.Logger.Debug(message)

	pars := readLog{
		Log: message,
		Time: utils.TimeFormat(),
		Level: "DEBUG",
	}

	ReadLog(pars)
}

func (logStruct *LogStruct) Warn(message string) {
	logStruct.Logger.Warn(message)

	pars := readLog{
		Log: message,
		Time: utils.TimeFormat(),
		Level: "WARN",
	}

	ReadLog(pars)
}

func (logStruct *LogStruct) Error(message string) {
	logStruct.Logger.Error(message)

	pars := readLog{
		Log: message,
		Time: utils.TimeFormat(),
		Level: "ERROR",
	}

	ReadLog(pars)
}

func ReadLog(pars readLog) {
	file, err := utils.FileEx(conf.Cfg.Path.LogPath)
	if err != nil {
		file, err = os.Create("../internal/utils/log/log_storage.json")
		if err != nil {
			Error.GetErr(err)
		}
	}

	defer file.Close()

	var logs []readLog
	data, err := io.ReadAll(file)
	if err != nil {
		Error.GetErr(err)
	}

	if len(data) > 0 {
		if unmarshalErr := json.Unmarshal(data, &logs); unmarshalErr != nil {
			Error.GetErr(unmarshalErr)
		}
	}

	logs = append([]readLog{pars}, logs...)

	utils.FileClear(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(logs)
	if err != nil {
		Error.GetErr(err)
	}
}


func LogInit() *LogStruct{
	var level slog.Level

	switch conf.Cfg.Env {
	case "local":
		level = slog.LevelDebug
	case "dev":
		level = slog.LevelInfo
	case "staging":
		level = slog.LevelWarn
	case "prod":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var log LogStruct
	log.Logger = *slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	return &log
}

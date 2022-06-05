package image_errors

import "errors"

var (
	ErrCantConnect   = errors.New("error: can't connect to Message broker")
	ErrChannelCreate = errors.New("error: can't create channel")
	ErrQueueCreate   = errors.New("error: can't create queue")
	ErrPublish       = errors.New("error: can't publish message")
	ErrConsume       = errors.New("error: can't consume message")
	ErrMarshal       = errors.New("error: can't marshal message")
	ErrUnmarshal     = errors.New("error: can't unmarshal message")
	ErrCantConvert   = errors.New("error: can't convert image")
	ErrFolderCreate  = errors.New("error: can't create folder")
	ErrFileCreate    = errors.New("error: can't create file")
	ErrSaveToFile    = errors.New("error: can't save image file")
	ErrFileOpen      = errors.New("error: can't open file")
	ErrImageDecode   = errors.New("error: can't decode image")
)

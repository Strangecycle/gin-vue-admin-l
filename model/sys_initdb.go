package model

type InitDBFunc interface {
	Init() error
}

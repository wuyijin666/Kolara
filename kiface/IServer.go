package kiface

type IServer interface {
	Start()
	Stop()
	Serve()
}	

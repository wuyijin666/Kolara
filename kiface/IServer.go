package kiface

type Server interface {
	Start()
	Stop()
	Serve()
}	

package structs

// Client Configuration struct
type ClientConfiguration struct {
	Socket_Path                string
	Log_Dir                    string
	Server_Log                 string
	Server_Connection_Attempts int
	Client_Log                 string
	Default_Command            string
	Default_Depth              int
	Default_FadeIn             int
	Default_FadeOut            int
	Default_Input              string
	Default_Loop               bool
	Default_Shuffle            bool
	Default_Volume             int
}

// Server Configuration struct
type ServerConfiguration struct {
	Socket_Path     string
	Log_Dir         string
	Server_Log      string
	Client_Log      string
	Default_Loop    bool
	Default_Shuffle bool
	Default_Volume  int
}

// Request struct
type Request struct {
	Command string
	Data    Data
}

// Data struct
type Data struct {
	Depth   int
	FadeIn  int
	FadeOut int
	Path    string
	Shuffle bool
	Loop    bool
	Values  []string
	Volume  int
}

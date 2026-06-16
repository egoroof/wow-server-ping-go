package ping

type server struct {
	Name string
	Host string
	Port int
}

type serverGroup struct {
	Name string
	List []server
}

var Servers = []serverGroup{
	{
		Name: "local",
		List: []server{
			{
				Name: "local",
				Host: "127.0.0.1",
				Port: 8085,
			},
		},
	},
	{
		Name: "x1",
		List: []server{
			{
				Name: "WoW Circle 3.3.5a x1",
				Host: "87.228.58.62",
				Port: 11294,
			},
			{
				Name: "WoW Circle 3.3.5a x1 [DE]",
				Host: "194.247.187.187",
				Port: 11294,
			},
			{
				Name: "WoW Circle 3.3.5a x1 [FIN]",
				Host: "193.84.2.209",
				Port: 11294,
			},
			{
				Name: "WoW Circle 3.3.5a x1 [NL]",
				Host: "31.207.45.133",
				Port: 11294,
			},
			{
				Name: "WoW Circle 3.3.5a x1 [MSK]",
				Host: "45.138.163.171",
				Port: 11294,
			},
			{
				Name: "WoW Circle 3.3.5a x1 [NSK]",
				Host: "79.141.77.15",
				Port: 11294,
			},
		},
	},
	{
		Name: "x4",
		List: []server{
			{
				Name: "WoW Circle 3.3.5a x4",
				Host: "178.72.132.84",
				Port: 10525,
			},
			{
				Name: "WoW Circle 3.3.5a x4 [DE]",
				Host: "194.247.187.136",
				Port: 10525,
			},
			{
				Name: "WoW Circle 3.3.5a x4 [FIN]",
				Host: "193.84.2.201",
				Port: 10525,
			},
			{
				Name: "WoW Circle 3.3.5a x4 [NL]",
				Host: "185.70.187.50",
				Port: 10525,
			},
			{
				Name: "WoW Circle 3.3.5a x4 [MSK]",
				Host: "5.35.8.118",
				Port: 10525,
			},
			{
				Name: "WoW Circle 3.3.5a x4 [NSK]",
				Host: "79.141.74.29",
				Port: 10525,
			},
		},
	},
	{
		Name: "x100",
		List: []server{
			{
				Name: "WoW Circle 3.3.5a x100",
				Host: "212.41.28.25",
				Port: 12742,
			},
			{
				Name: "WoW Circle 3.3.5a x100 [DE]",
				Host: "194.247.187.187",
				Port: 12742,
			},
			{
				Name: "WoW Circle 3.3.5a x100 [FIN]",
				Host: "193.84.2.209",
				Port: 12742,
			},
			{
				Name: "WoW Circle 3.3.5a x100 [NL]",
				Host: "31.207.45.133",
				Port: 12742,
			},
			{
				Name: "WoW Circle 3.3.5a x100 [MSK]",
				Host: "45.138.163.171",
				Port: 12742,
			},
			{
				Name: "WoW Circle 3.3.5a x100 [NSK]",
				Host: "79.141.77.15",
				Port: 12742,
			},
		},
	},
	{
		Name: "Fun",
		List: []server{
			{
				Name: "WoW Circle 3.3.5a Fun",
				Host: "87.228.3.124",
				Port: 12373,
			},
			{
				Name: "WoW Circle 3.3.5a Fun [DE]",
				Host: "194.247.187.187",
				Port: 12373,
			},
			{
				Name: "WoW Circle 3.3.5a Fun [FIN]",
				Host: "193.84.2.209",
				Port: 12373,
			},
			{
				Name: "WoW Circle 3.3.5a Fun [NL]",
				Host: "31.207.45.133",
				Port: 12373,
			},
			{
				Name: "WoW Circle 3.3.5a Fun [MSK]",
				Host: "45.138.163.171",
				Port: 12373,
			},
			{
				Name: "WoW Circle 3.3.5a Fun [NSK]",
				Host: "79.141.77.15",
				Port: 12373,
			},
		},
	},
}

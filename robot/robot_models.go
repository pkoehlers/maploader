package robot

var robots []Robot = []Robot{
	/*
		{model: "your.hostname", // Debug config
			MapFolders:       []string{},
			mapFiles:         []string{},
			restartProcesses: []Process{}},
	*/
	{model: "p2029", // Dreame L10 Pro
		MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap"},
		mapFiles:         []string{"/data/config/ava/mult_map.json"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "p2008", // Dreame F9
		MapFolders:       []string{"/data/log/ri/", "/data/DivideMap"},
		mapFiles:         []string{"/data/config/ava/mult_map.json", "/data/log/map_info.bin", "/data/log/slam.db"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "p2009", // Dreame D9
		MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap", "/data/DivideDebug"},
		mapFiles:         []string{"/data/config/ava/mult_map.json", "/data/log/map_info.bin"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "p2187", // Dreame D9 Pro
		MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap", "/data/DivideDebug"},
		mapFiles:         []string{"/data/config/ava/mult_map.json", "/data/log/map_info.bin"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "p2027", // Dreame W10
		MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap", "/data/DivideDebug"},
		mapFiles:         []string{"/data/config/ava/mult_map.json", "/data/log/map_info.bin"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "r2228", // Dreame L10s Ultra
		MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap"},
		mapFiles:         []string{"/data/config/ava/mult_map.json"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "r2240", // Dreame D10S Plus
		MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap"},
		mapFiles:         []string{"/data/config/ava/mult_map.json"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "r2250", // Dreame D10S Pro
		MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap"},
		mapFiles:         []string{"/data/config/ava/mult_map.json"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
	{model: "r2416", // Dreame X40
		MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap", "/data/DivideDebug"},
		mapFiles:         []string{"/data/config/ava/mult_map.json", "/data/log/map_info.bin"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
}

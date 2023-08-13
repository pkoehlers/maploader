package robot

var robots []Robot = []Robot{
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
}

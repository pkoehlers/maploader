package robot

var robots []Robot = []Robot{
	{model: "p2029", // Dreame L10 Pro
		MapFolders:       []string{"/data/ri", "/data/map", "/data/DivideMap"},
		mapFiles:         []string{"/data/config/ava/mult_map.json"},
		restartProcesses: []Process{MiioClientProcess, AvaProcess}},
}

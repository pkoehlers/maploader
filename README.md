# Dreame Vacuum Robot Maploader

Provides a map changing functionality for rooted vacuum robots running Valetudo, controllable via MQTT. This project currently supports Dreame robots, see the supported model list below.
A similar project for Xiaomi/Roborock vacuums (without further affiliation) can be found here: [Thyraz/MapLoader](https://github.com/Thyraz/MapLoader).

## Note
The map changing process used in this project is based on observations and testing. It is not officially supported by Dreame or Valetudo. It's recommended to back up your map before using the maploader. Map changes should not be done during cleaning tasks but may be done when the robot is not docked.

## How it works
The maploader is a program running on the robot. It communicates with a configured MQTT broker to manage map changes. 

When a new map name is received, the robot backs up the current map, removes all map files, and loads the new map if it exists. The robot's relevant processes are then restarted to apply the new map.

In case of issues, backup archives for each map are available in `/data/maploader`. 

The default map name is "main". 

After map change, Valetudo restarts and may temporarily show an empty or the old map. This process can be sped up by starting the cleaning (and stopping it immediately).

## Feature: WAV Audio Notification
Maploader now supports playing a WAV audio file for status changes. This feature is particularly useful for audible notifications when the map is loaded or changed. The audio file path and arguments for the `aplay` command are configurable via environment variables.

### Configuration for WAV Audio Feature
- **WAV_FILE_MAP_LOADED**: Set the path to the WAV file. You can use `{map_name}` as a placeholder in the path, which will be replaced with the current map name. If the specific map's WAV file does not exist, a default file (with "default" replacing `{map_name}`) will be used.
- **WAV_APLAY_ARGS**: Optionally set the arguments for the `aplay` command. Defaults to `-Dhw:0,0` if not set.

Below is an example of the filenames for the WAV files, showcasing different audio files for different maps:
- `/data/maploader/wav`
  - `map_change_default.wav`
  - `map_change_livingroom.wav`
  - `map_change_main.wav`

This can be used with this environment variable:
`WAV_FILE_MAP_LOADED=/data/maploader/wav/map_change_{map_name}.wav`

The wav directory includes a default file (`map_change_default.wav`) and specific files for different map names (e.g., `livingroom`, `main`).


The WAV file will be played automatically after map changes to provide auditory feedback.

## MQTT Topics

| Topic                                      | Publisher    | Payload   | Description                                                                                                                                                                                                          |
|--------------------------------------------|--------------|-----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `valetudo/{identifier}/maploader/map`      | Robot + Home | Map name  | Maploader publishes the current map name to this topic after a map switch.<br/><br/>Can be used to set the current map name without changing the active map.                                                         |
| `valetudo/{identifier}/maploader/map/set`  | Home         | Map name  | Stores the active map and switches to the given map.<br/><br/>A backup of the current map file is made before the active map is stored.<br/><br/>If no map is stored under the given map name a blank map is loaded. |
| `valetudo/{identifier}/maploader/map/save` | Home         | Map name  | Saves the active map under the given name and switches to that map.<br/><br/>A backup of the map file is made before it is overwritten.                                                                              |
| `valetudo/{identifier}/maploader/map/load` | Home         | Map name  | Switches to the given map without storing the active map.<br/><br/>If no map is stored under the given map name a blank map is loaded.                                                                               |
| `valetudo/{identifier}/maploader/status`   | Robot        | see below | Maploader publishes its current status to this topic.                                                                                                                                                                |

The identifier is set in the MQTT settings in Valetudo.

The maploader status can change to the following value:

| Status            | Description                                             |
|-------------------|---------------------------------------------------------|
| idle              | Maploader is started and awaiting commands              |
| loading_map       | The map is currently being loaded                       |
| saving_map        | The map is currently being saved                        |
| starting_services | Waiting for robot services to restart                   |
| error             | An error occurred, logs need to be checked              |
| offline           | The maploader process exited / lost the MQTT connection |

## Home Assistant Config
This project does not support Home Assistant auto discovery as I am using the sensor to define the list of possible maps. To allow Home Assistant to work with maploader add the section below to your configuration.yaml. To create a new map, just add a new value to the field and set the entity to that new value.

```
mqtt:
  sensor:
    - state_topic: valetudo/foo/maploader/status
      name: "vacuum_maploader_status"
  select:
    - command_topic: valetudo/foo/maploader/map/set
      state_topic: valetudo/foo/maploader/map
      name: "vacuum_maploader_map"
      options:
        - "main"
        - "second_floor"

```

# Supported robots

The following models are known to work with the maploader:

| Model                | binary                     |
|----------------------|----------------------------|
| Dreame L10 Pro       | maploader-arm64            |
| Dreame Z10 Pro       | maploader-arm64            |
| Dreame L10s Ultra    | maploader-arm64            |
| Dreame F9            | maploader-arm              |
| Dreame D9            | maploader-arm              |
| Dreame D9 Pro        | maploader-arm              |
| Dreame W10           | maploader-arm              |

If your Dreame robot is not listed here, you need to find out the arch of your robot (e.g. with ```uname -m```, where ```aarch64``` -> ...-arm64 binary).

Please open an issue, if your Dreame robot is not listed here but works or you need assistance.

# Installation
The binary must be placed in the ```/data``` folder and it needs to be started with the system.

Make sure that Valetudo is working and MQTT is also setup in Valetudo.

## Steps

Open a SSH terminal

Download the binary: (check the supported robots section for the correct filename in the URL)

```wget -O /data/maploader-binary https://github.com/pkoehlers/maploader/releases/latest/download/maploader-arm64```

Add execution permissions:

```chmod +x /data/maploader-binary```

Open the postboot script:

```vi /data/_root_postboot.sh```

Search for the valetudo start block and add the maploader binary after the start of Valetudo so that it looks like this:
```
if [[ -f /data/valetudo ]]; then
        VALETUDO_CONFIG_PATH=/data/valetudo_config.json /data/valetudo > /dev/null 2>&1 &
        VALETUDO_CONFIG_PATH=/data/valetudo_config.json /data/maploader-binary > /dev/null 2>&1 &
fi
```

Reboot

# Configuration
The following things can be customized with environment variables:

| Variable                   | Default value              | Description                                                    |
|----------------------------|----------------------------|----------------------------------------------------------------|
| VALETUDO_CONFIG_PATH       | /data/valetudo_config.json | Path to the valetudo config file                               |
| MAPLOADER_RESTART_VALETUDO |                            | Set to any value to restart Valetudo after map changes         |
| MAPLOADER_DIR              | /data/maploader            | Directory to store the maps and logs                           |
| DEFAULT_MAP_NAME           | main                       | Name of the main map (for first use)                           |
| ROTATION_KEEP_MAPS         | 5                          | Number of map backups to keep per map                          |
| WAV_FILE_MAP_LOADED        |                            | Path to WAV file for audio notification (supports placeholders)|
| WAV_APLAY_ARGS             | -Dhw:0,0                   | Arguments for the `aplay` command (audio playback)             |


# Technical Details
As mentioned this is only tested with the Dreame L10 Pro but other Dreame robots should work just fine.
Currently, these files/directories are considered "map files":

* ```/data/ri```
* ```/data/map```
* ```/data/DivideMap```
* ```/data/config/ava/mult_map.json```

Basically these are all the files that will be cleared by resetting the map via Valetudo.

# Development
Build the project with
```env GOOS=linux GOARCH=arm64 go build -o maploader-binary```

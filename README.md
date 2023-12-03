# Dreame Vacuum Robot Maploader

Provides a map changing functionality for rooted vacuum robots running Valetudo controllable via MQTT.

Currently this only supports Dreame robots, tested with a [Dreame L10 Pro](https://dontvacuum.me/robotinfo/detail_dreame.vacuum.p2029_0.html).

A similar project for Xiaomi/Roborock vacuums (without further affiliation) can be found here: [Thyraz/MapLoader](https://github.com/Thyraz/MapLoader)

## Note
The map changing process used in this project is purely based on observations and testing. 
This is neither supported by Dreame nor Valetudo. If you consider your map valueable, back it up before using the maploader.

It does not matter if the vacuum is docked or not but ensure that map changes don't take place during cleaning tasks.
# How it works
The maploader is a small programm running on the robot. It creates a state and command topic in the configured MQTT broker. 

When a new map name is sent to the command topic, the robot will backup the current map files, remove all the map files and restore the other map if it exists. It will finally restart relevant processes on the vacuum to load the new map.

In case anything goes wrong, you can get the last few backup archives for each map in ```/data/maploader```.

The default map name is "main".

After the map change, Valetudo will be restarted and will not be reachable for some seconds. Sometimes Valetudo might show an empty map after restarting and it takes some time to load the actual map. This process can be sped up by starting the cleaning (and stopping it directly).

I am using this with Home Assistant, where I trigger the map change as part of an automation and move the robot to the other zone. It can then be operated on the new map after the reboot.

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

| Variable                   | Default value              | Description                                                 |
|----------------------------|----------------------------|-------------------------------------------------------------|
| VALETUDO_CONFIG_PATH       | /data/valetudo_config.json | path to the valetudo config file                            |
| MAPLOADER_RESTART_VALETUDO |                            | Set to any value to also restart Valetudo after map changes |
| MAPLOADER_DIR              | /data/maploader            | directory to store the maps and logs                        |
| DEFAULT_MAP_NAME           | main                       | name of the main map (for first use)                        |
| ROTATION_KEEP_MAPS         | 5                          | number of map backups to keep per map                       |



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

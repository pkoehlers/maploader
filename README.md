# Dreame Vacuum Robot Maploader

Provides a map changing functionality for rooted vacuum robots running Valetudo controllable via MQTT.

Currently this only supports Dreame robots, tested with a [Dreame L10 Pro](https://dontvacuum.me/robotinfo/detail_dreame.vacuum.p2029_0.html).

A similar project for Xiaomi/Roborock vacuums (without further affiliation) can be found here: [Thyraz/MapLoader](https://github.com/Thyraz/MapLoader)

# How it works
The maploader is a small programm running on the robot. It creates a state and command topic in the configured MQTT broker. The default map name is "main". When another name is sent to the command topic, the robot will backup the current map files, remove all the map files and restore the other map if it exists. It will finally reboot the robot.

I am using this with Homeassistant, where I trigger the map change as part of an automation and move the robot to the other zone. It can then be operated on the new map after the reboot.

# Installation
The binary must be placed in the ```/data``` folder and it needs to be started with the system.

Make sure that Valetudo is working and MQTT is also setup in Valetudo.

## Steps

Open a SSH terminal

Download the binary:

```wget -O /data/maploader-binary https://github.com/pkoehlers/maploader/releases/download/v1.0.0/maploader-arm64```

Add execution permissions:

```chmod +x /data/maploader-binary```

Open the postboot script:

```vi /data/_root_postboot.sh```

Search for the valetudo start block and add the maploader statup after the start of Valetudo so that it looks like this:
```
if [[ -f /data/valetudo ]]; then
        VALETUDO_CONFIG_PATH=/data/valetudo_config.json /data/valetudo > /dev/null 2>&1 &
        VALETUDO_CONFIG_PATH=/data/valetudo_config.json /data/maploader-binary > /dev/null 2>&1 &
fi
```

Reboot


# Technical Details
As mentioned this is only tested with the Dreame L10 Pro but other Dreame robots should work just fine.
Currently these files/direcotries are considered "map files":

* ```/data/ri```
* ```/data/map```
* ```/data/DivideMap```
* ```/data/config/ava/mult_map.json```

Basically these are all the files that will be cleared by resetting the map via Valetudo.

# pumpkin-pi 🎃

![Build Status](https://github.com/circa10a/pumpkin-pi/workflows/build-docker-image/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/circa10a/pumpkin-pi)](https://pkg.go.dev/github.com/circa10a/pumpkin-pi?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/github.com/circa10a/pumpkin-pi)](https://goreportcard.com/report/github.com/circa10a/pumpkin-pi)

Raspberry pi project that controls jack-o-lantern via servo motor and PIR motion sensors to simulate it "watching" you.

Inspired by [Ryder Damen's mannequin head](https://www.youtube.com/watch?v=9CVhrZEhEoE)

- [Demo](#demo)
- [Photos](#photos)
- [Wiring diagram](#wiring-diagram)
- [Deploy](#deploy)
  - [Docker](#docker)
  - [Go](#Go)
- [Materials](#materials)
- [Configuration](#configuration)

## Demo

> This project was originally built for outdoor use

[Youtube link](http://www.youtube.com/watch?v=fl52GQJCFVI)

## Photos

> Note: Holes in the acrylic case are needed for the motion sensors to properly work

<p float="left">
  <img src="https://i.imgur.com/ngqizRO.jpg" width="45%" height="45%"/>
  <img src="https://i.imgur.com/uXo3kaP.jpg" width="45%" height="45%"/>
  <img src="https://i.imgur.com/IPJ1QB8.jpg" width="45%" height="45%"/>
  <img src="https://i.imgur.com/YtHwhJn.jpg" width="45%" height="45%"/>
  <img src="https://i.imgur.com/9AHhIF6.jpg" width="45%" height="45%"/>
  <img src="https://i.imgur.com/snKIZKp.jpg" width="45%" height="45%"/>
<p/>

## Wiring diagram

> Created with [circuit-diagram.org](https://www.circuit-diagram.org/). Source file is in [/diagrams](/diagrams)

![alt text](/images/circuit.png)

## Deploy

> Requires following diagram above to be wired up correctly

### Docker

```bash
# This script will install the dependencies and start the containers
bash -c "$(curl -sL https://raw.githubusercontent.com/circa10a/pumpkin-pi/main/install.sh)"
```

### Go

> Requires Go 1.17+

1. Install [Go](https://golang.org/doc/install)
2. Install [pi-blaster](https://github.com/sarfata/pi-blaster)
3. `go install github.com/circa10a/pumpkin-pi@latest`
4. `pumpkin-pi`

## Materials

- [Raspberry Pi (This project uses a Pi 3 Model B)](https://www.adafruit.com/product/4292)
- [USB Power Supply Compatible with Pi3/4](https://www.amazon.com/dp/B07X8C6PV6/ref=cm_sw_em_r_mt_dp_02B2MZVAR88S0RR0J65M)
- [MicroSD card](https://www.amazon.com/dp/B004KSMXVM/ref=cm_sw_em_r_mt_dp_7JES2YT0FC79MHFDBP1Z) with Raspbian installed ([guide](https://www.raspberrypi.org/documentation/computers/getting-started.html#installing-the-operating-system))
- [Pimoroni pan-tilt hat](https://shop.pimoroni.com/products/pan-tilt-hat?variant=22408353287)
- [2x PIR motion sensors](https://www.amazon.com/dp/B07KBWVJMP/ref=cm_sw_em_r_mt_dp_JHKZXKE9W8X144C21QZX)
- [2x 3D printed PIR motion sensor enclosures](https://www.thingiverse.com/thing:3366814)
- [MG90s Servo Motor](https://www.amazon.com/dp/B07NV476P7/ref=cm_sw_em_r_mt_dp_HMWNXQVMQZKX0D4K25SD)
  - I upgraded the horizontal motor in the pan-tilt hat due to needing more torque to support/smoothly move the pumpkin head without struggle
- [Male to female, female to femaile jumper wires](https://www.amazon.com/dp/B01EV70C78/ref=cm_sw_em_r_mt_dp_SWRTQ805V399FCG4DCFH)
- [M2 Standoffs](https://www.amazon.com/dp/B07B9X1KY6/ref=cm_sw_em_r_mt_dp_WZWF9MSF0CDSYY296XG6)
- [Jack-o-lantern (with top stem cut off)](https://www.homedepot.com/p/Home-Accents-Holiday-9-in-White-Blow-Mold-Pumpkin-with-Black-Shadow-21GM27288/315532374)
- [Acrylic display case](https://www.hobbylobby.com/search/?text=display+case&quickview=81011632)
- [Plaster column pedastal(spray painted blacked)](https://www.hobbylobby.com/Home-Decor-Frames/Furniture/Accent-Furniture/White-Corinthian-Column-Pedestal/p/CP02000)
- [2x 5v fans for (intake/exhaust)](https://www.amazon.com/dp/B07KRSJVP7/ref=cm_sw_em_r_mt_dp_G485R3B54ETDE8D3KKQZ)
- Some screws here and there

## Configuration

|                                               |                                                                                    |           |          |
|-----------------------------------------------|------------------------------------------------------------------------------------|-----------|----------|
| Environment Variable                          | Description                                                                        | Required  | Default  |
| `PUMPKINPI_LOG_LEVEL`                         | [Logrus](https://github.com/sirupsen/logrus) log level. Use `debug` for more info  | `false`   | `info`   |
| `PUMPKINPI_MOTION_TIMES_ENABLED`              | Whether to use configured schedule or not. These times must be within the same day | `false`   | `false`  |
| `PUMPKINPI_MOTION_TIME_START`                 | Local time to ensure pumpkin-pi only responds after this hour                      | `false`   | `17`     |
| `PUMPKINPI_MOTION_TIME_END`                   | Local time to ensure pumpkin-pi only responds before this hour                     | `false`   | `22`     |
| `PUMPKINPI_PIR_LEFT_MOTION_SENSOR_GPIO_PIN`   | The GPIO Pin used to read inputs from left motion sensor                           | `false`   | `11`     |
| `PUMPKINPI_PIR_RIGHT_MOTION_SENSOR_GPIO_PIN`  | The GPIO Pin used to read inputs from right motion sensor                          | `false`   | `13`     |
| `PUMPKINPI_SERVO_CENTER`                      | The center position of the horizontal servo motor                                  | `false`   | `32`     |
| `PUMPKINPI_SERVO_LEFT`                        | The left position of the horizontal servo motor                                    | `false`   | `23`     |
| `PUMPKINPI_SERVO_RIGHT`                       | The right position of the horizontal servo motor                                   | `false`   | `40`     |
| `PUMPKINPI_SERVO_ROTATE_DELAY`                | The wait time in between each incremental servo step as it rotates                 | `false`   | `150ms`  |
| `PUMPKINPI_SERVO_CENTER_RESET_INTERVAL`       | The interval at which the pumpkin will rotate back to the center position          | `false`   | `5m`     |
| `PUMPKINPI_SERVO_GPIO_PIN`                    | The PWM enabled GPIO Pin used to control the servo motor                           | `false`   | `12`     |

package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	// All environment variables for configuration expect PUMPKINPI_ prefix
	envVarPrefix = "pumpkinpi"
)

// PumpkinPiConfig is a global config struct which holds all configuration options
type PumpkinPiConfig struct {
	LogLevel           string `envconfig:"LOG_LEVEL" default:"debug"`
	MotionTimesEnabled bool   `envconfig:"MOTION_TIMES_ENABLED" default:"false"`
	MotionTimeEnd      int    `envconfig:"MOTION_TIME_END" default:"22"`
	MotionTimeStart    int    `envconfig:"MOTION_TIME_START" default:"18"`
	// Ensure multiple events doesn't jerk motor back and forth with this cheap "lock"
	MovingLock                  bool
	ServoCenter                 uint8         `envconfig:"SERVO_CENTER" default:"32"`
	ServoCenterResetInterval    time.Duration `envconfig:"SERVO_CENTER_RESET_INTERVAL" default:"5m"`
	ServoLeft                   uint8         `envconfig:"SERVO_LEFT" default:"23"`
	ServoRight                  uint8         `envconfig:"SERVO_RIGHT" default:"40"`
	ServoRotateDelay            time.Duration `envconfig:"SERVO_ROTATE_DELAY" default:"150ms"`
	ServoGPIOPin                string        `envconfig:"SERVO_GPIO_PIN" default:"12"`
	PIRLeftMotionSensorGPIOPin  string        `envconfig:"PIR_LEFT_MOTION_SENSOR_GPIO_PIN" default:"11"`
	PIRRightMotionSensorGPIOPin string        `envconfig:"PIR_RIGHT_MOTION_SENSOR_GPIO_PIN" default:"13"`
}

// isDuringConfiguredHours will determine if the pumpkin-pi should physically respond if during configured motion times
func (p *PumpkinPiConfig) isDuringConfiguredHours(currentHour, startHour, endHour int) bool {
	if p.MotionTimesEnabled {
		return currentHour >= startHour && currentHour < endHour
	}
	return true
}

func (p *PumpkinPiConfig) initLogger() (*log.Logger, error) {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	level, err := log.ParseLevel(p.LogLevel)
	if err != nil {
		return log.New(), err
	}
	logger.SetLevel(level)
	return logger, nil
}

func main() {
	// Init configuration options
	pumpkinPiConfig := &PumpkinPiConfig{}
	err := envconfig.Process(envVarPrefix, pumpkinPiConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Logger Config
	log, err := pumpkinPiConfig.initLogger()
	if err != nil {
		log.Fatal(err)
	}
	// Current time for debugging
	log.Debug("current time: ", time.Now())
	// Init adapter
	r := raspi.NewAdaptor()
	// Configure drivers for hardware devices
	servo := gpio.NewServoDriver(r, pumpkinPiConfig.ServoGPIOPin)
	leftSensor := gpio.NewPIRMotionDriver(r, pumpkinPiConfig.PIRLeftMotionSensorGPIOPin)
	rightSensor := gpio.NewPIRMotionDriver(r, pumpkinPiConfig.PIRRightMotionSensorGPIOPin)
	// Set servo back to center position
	log.Debug("setting servo to center position")
	err = servo.Move(pumpkinPiConfig.ServoCenter)
	if err != nil {
		log.Error(err)
	}
	// Track position for rotating back to center
	currentPosition := pumpkinPiConfig.ServoCenter

	// the meat
	work := func() {
		// Left motion sensor
		err = leftSensor.On(gpio.MotionDetected, func(data interface{}) {
			log.Debug("left motion detected")
			// If during configured hours
			if pumpkinPiConfig.isDuringConfiguredHours(time.Now().Hour(), pumpkinPiConfig.MotionTimeStart, pumpkinPiConfig.MotionTimeEnd) {
				log.Debug("current time is between motion times or motion times are disabled")
				// Ensure pumpkin is not already moving from another event
				if pumpkinPiConfig.MovingLock {
					return
				}
				if !pumpkinPiConfig.MovingLock {
					pumpkinPiConfig.MovingLock = true
					// Rotate motor left incrementally
					log.Debug("setting servo to left position")
					for i := currentPosition; i >= pumpkinPiConfig.ServoLeft; i-- {
						time.Sleep(pumpkinPiConfig.ServoRotateDelay)
						err = servo.Move(i)
						if err != nil {
							log.Error(err)
						}
						currentPosition = pumpkinPiConfig.ServoLeft
					}
					pumpkinPiConfig.MovingLock = false
				}
			}
		})
		if err != nil {
			log.Error(err)
		}
		// Right motion sensor
		err = rightSensor.On(gpio.MotionDetected, func(data interface{}) {
			log.Debug("right motion detected")
			// If during configured hours
			if pumpkinPiConfig.isDuringConfiguredHours(time.Now().Hour(), pumpkinPiConfig.MotionTimeStart, pumpkinPiConfig.MotionTimeEnd) {
				log.Debug("current time is between motion times or motion times are disabled")
				// Ensure pumpkin is not already moving from another event
				if pumpkinPiConfig.MovingLock {
					return
				}
				if !pumpkinPiConfig.MovingLock {
					pumpkinPiConfig.MovingLock = true
					// Rotate motor right incrementally
					log.Debug("setting servo to right position")
					for i := currentPosition; i <= pumpkinPiConfig.ServoRight; i++ {
						time.Sleep(pumpkinPiConfig.ServoRotateDelay)
						err = servo.Move(i)
						if err != nil {
							log.Error(err)
						}
						currentPosition = pumpkinPiConfig.ServoRight
					}
					pumpkinPiConfig.MovingLock = false
				}
			}
		})
		if err != nil {
			log.Error(err)
		}
		// Reset motor back to center periodically
		gobot.Every(pumpkinPiConfig.ServoCenterResetInterval, func() {
			log.Debug("executing reset back to center scheduler function")
			// If during configured hours
			if pumpkinPiConfig.isDuringConfiguredHours(time.Now().Hour(), pumpkinPiConfig.MotionTimeStart, pumpkinPiConfig.MotionTimeEnd) {
				log.Warn(time.Now().Hour(), pumpkinPiConfig.MotionTimeStart, pumpkinPiConfig.MotionTimeEnd)
				log.Debug("current time is between motion times or motion times are disabled")
				// Ensure pumpkin is not already moving from another event
				if !pumpkinPiConfig.MovingLock {
					pumpkinPiConfig.MovingLock = true
					// If motor is in the right position
					if currentPosition > pumpkinPiConfig.ServoCenter {
						// Rotate motor left incrementally
						for i := currentPosition; i >= pumpkinPiConfig.ServoCenter; i-- {
							log.Debug("pumpkin currently set to right position. setting servo back to center position due to scheduler")
							time.Sleep(pumpkinPiConfig.ServoRotateDelay)
							err = servo.Move(i)
							if err != nil {
								log.Error(err)
							}
						}
					}
					// If motor is in the left position
					if currentPosition < pumpkinPiConfig.ServoCenter {
						// Rotate motor right incrementally
						for i := currentPosition; i <= pumpkinPiConfig.ServoCenter; i++ {
							log.Debug("pumpkin currently set to left position. setting servo back to center position due to scheduler")
							time.Sleep(pumpkinPiConfig.ServoRotateDelay)
							err = servo.Move(i)
							if err != nil {
								log.Error(err)
							}
						}
					}
					currentPosition = pumpkinPiConfig.ServoCenter
					pumpkinPiConfig.MovingLock = false
				}
			}
		})
	}
	robot := gobot.NewRobot(
		"pumpkin-pi",
		[]gobot.Connection{r},
		[]gobot.Device{leftSensor, rightSensor, servo},
		work,
	)
	err = robot.Start()
	if err != nil {
		log.Fatal(err)
	}
}

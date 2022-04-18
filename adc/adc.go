// @file TinyGo/adc/adc.go
// @date 2022.02.28 fork from TinyGo/adc_pi/adc/adc.go --warning-- same pkg name
// @info capture buffer of adc samples; audio capture rate

// @require: go mod init localhost/adc
// @build: tinygo flash -target=pico

// @date 2022.03.01 added %04x formating to Cap2Uart output; fixes uart_xfr char count bug
// @date 2022.03.02 added Cap2UartThreshold, ignores sound below threshold
// @date 2022.03.10 added Cap2Uint16, updated Cap2Uart to call it
// @date 2022.03.24 changed adc_cap_threshold* for live capture
// @date 2022.03.25 removed zero padding for capture < 1024
// @date 2022.04.08 re-enabled lastSoundPos; disabled threshold_low
// @date 2022.04.09 changed threshold from 2.0V to 1.75V
// @date 2022.04.18 removed import 'common'; const Tag* added locally

package adc

import (
        "fmt"
	"machine"
	"time"
)

const adc_cap_threshold     = 35000 // 1.75V (/ (* 1.75 65536) 3.3) 34753
// --obs-- const adc_cap_threshold_low = 20000 // 1.0V (/ (* 1.0 65536) 3.3) 19859

// general purpose tags; copied from 'common' and removed the localhost/common dependency
const Tag_file     = "--file--"
const Tag_eod      = "--eod--" // end of data
const Tag_eot      = "--eot--" // end of transmission
const Out_file     = "not-in-git.txt" // scratch file, e.g. created by dsp.Pull()

// Cap2Uart captures 'buf_size' samples from adc with sample time of 'sleep_time' + Get() us
func Cap2Uart(buf_size, sleep_time int) {
	tag_file := Tag_file
	tag_eod  := Tag_eod
	tag_eot  := Tag_eot

	buf := Cap2Uint16( buf_size, sleep_time )
	
	fmt.Printf("........%s--\n\r", tag_file)
	for _,v := range buf {
		fmt.Printf("%04x\n\r", v)
	}
	fmt.Printf("%s\n\r", tag_eod)
	fmt.Printf("%s\n\r", tag_eot)
} // end func Cap2Uart2(buf_size, sleep_time int) 

// Cap2Uint16 captures, processes, and returns adc data; const local adc.go threshold values
func Cap2Uint16(buf_size, sleep_us int) (buf []uint16){
	threshold := adc_cap_threshold // const atop adc.go
	// --obs-- threshold_low := adc_cap_threshold_low // const atop adc.go
	machine.InitADC()
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	sensor := machine.ADC{machine.ADC0}
	sensor.Configure(machine.ADCConfig{})
	// --obs-- assume caller handles ui: fmt.Printf("Tinygo/adc Cap2Uint16 --blocking--\n\r")
	buf = make([]uint16, buf_size) // capture  buffer
	val := sensor.Get() // uint16 disposable first adc read initializes val
	led.High() // high when adc is blocking for threshold
	for { // wait for adc to exceed threshold
		val = sensor.Get() // uint16
		// sound input threshold;
		if val > uint16(threshold) {
			break;
		}
		buf[0] = val // first sample excluded from range below
	} // end wait for adc to exceed threshold
	led.Low()
	// --CAPTURE--
	for i:=1; i<len(buf); i++ { // range buf adc get; already have buf[0]
		// 'Get()' takes ~16us on pico?; 70 us sleep -> 86 us/samp
		// (+ 70 16) 86 (/ 1.0 86e-6) 11.6 Ksamp/sec
		// (+ 300 16) 316 (/ 1.0 316e-6) 3.16 Ksamp/sec
		buf[i] = sensor.Get() // uint16
		time.Sleep(time.Microsecond * time.Duration(sleep_us))
		// buf_size=2048, sleep_time=300 -> (* 316 2048 ) ~ 0.647168 second recording
	} // end range buf
	// end --CAPTURE--
	// fmt.Println("--debug-- buf[i]", buf[0:32], "\n\r")

	lastSoundPos := len(buf)-1 // find end of sound over threshold, and prune
	for i:=len(buf)-1; i>=0; i-- {
		// --obs-- if buf[i] >= uint16(threshold) || buf[i] <= uint16(threshold_low) {
		if buf[i] >= uint16(threshold) {
			lastSoundPos = i
			break
		}
	}

	// lastSoundPos = len(buf)-1 // --dev-- 20220408 disables lastSoundPos

	return buf[:lastSoundPos]
} // end func Cap2Uint16

// Notes:
//
// @info parse ascii hex to int:
//	threshold, e_th := strconv.ParseInt(Threshold, 0, 64) // hex to int
//	if err_th != nil {
//		fmt.Fprintf(os.Stderr, "--Error-- Cap2Uart() ParseInt()\n")
//		panic(e_th)
//	}
//



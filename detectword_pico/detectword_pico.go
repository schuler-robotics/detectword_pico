// @file TinyGo/detectword_pico/detectword_pico.go 
// @date 2022.03.08
// @info detect word pairs, e.g. 'on' and 'off', 'light' and 'dark' and set gpio signifiers accordingly
// @info requires "\n\r" newlines for tinygo/pico apps unless otherwise noted

// @build: tinygo flash -target=pico

// Copyright 2022 RC Schuler. All rights reserved.
// Use of this source code is governed by a GNU V3
// license that can be found in the LICENSE file.

// @date 2022.03.08 fork of sandbox_dw.go; first working demo of detectword_pico git tagged as 'demo_v0.1';
//                  threshold, threshold_low from adc/Cap2Array() come from Go/common
// @date 2022.03.11 working and matched to sandbox_dw on raspi; both use adc.Cap2Uint16, buf_size 2048 samples with
//                  possible 1024 zero padding, sleep_time 300us, Fbins = 64, Tbins == 64
// @date 2022.03.12 added gpio10 to control light
// @date 2022.03.15 added ReduceWordDetect()
// @date 2022.03.16 added ReduceWordDetectCreateRef(); Ref calls outside of loop
// @date 2022.03.24 added loopCt; set first two spoken words as light/dark ref samples; include file vars
//                  still loaded and unused.
// @date 2022.03.25 cleaned and minimized stdio
// @date 2022.03.27 added CreateU16SpectFromU16, moved Hamming() and complexFloatArray allocation outside func
// @date 2022.03.31 added MinWordLen (samples) to detection loop; added (unused) lightState; MinWordLen appears
//                  to have no effect, even set at 0.95 with buf_size=2048
// @date 2022.04.01 checking sound threshold with Normalize_ac_threshold(), 0xBFFF is 0.75 0.FFFF;
//                  checking sould levels on ref words before accepting
// @date 2022.04.08 cleanup; memory reduce low hanging fruit; starting point 20220403-working-restore
// @date 2022.04.09 detection tuning: changed buf_size, sleep_time from 2048, 500, to 1024, 100
// @date 2022.04.11 added --uart out-- calls uartHeader(), uartFooter(); tagging params from common;
//                  output xt, spect, pool1, and pool2 to uart for raspi uart_xfr acquisition
// @date 2022.04.13 commented Print*, capture_diags false for --prod--
//                  --prod-- Tbins/Fbins to 32 from 64, SpectThresh to 60 from 50
// @date 2022.04.18 removed import 'common'; const Tag* added locally

package main

import (
	"fmt" // --quiet-- mode
	"time"
	// "runtime" // runtime.GC is disabled
	"localhost/adc"     // underscore disable for --no mic-- mode
	"machine"
)

// general purpose tags; copied from 'common' and removed the localhost/common dependency
const Tag_file     = "--file--"
const Tag_eod      = "--eod--" // end of data
const Tag_eot      = "--eot--" // end of transmission
const Out_file     = "not-in-git.txt" // scratch file, e.g. created by dsp.Pull()

// Acquire first 2 data sets (words) as references for subsequent captures;  e.g. "on" and "off"
func main() {
	// --quiet-- fmt.Printf("\n\r## detectword_pico %s\n\r", fmt.Sprintf("%s",time.Now())[:16])
	time.Sleep(time.Millisecond * 1000) // power stabalize; added 20220401; usb batt #1 producing connect bounce
	
	// adc and spectrograph parameters
	Tbins := 64 // --prod-- 64
	Fbins := 64 // --prod-- 64
	buf_size := 1024  // --prod-- 1024 
	sleep_time := 250 // --prod-- 250; 'sleep_time'us + 16us == 'adc.Get' time; 'Tsamp' in octave  mfiles
	SpectThresh := uint16(50)  // --prod-- 50; ignore spect array elements below SpectThresh
	MinWordLen := int(0.2 * float64(buf_size)) // don't process sounds less than X% of buf_sizes
	// Reduction params
	vBlocks  := 8;  hBlocks := 8  // reduction block size for avg pool; require power of 2
	vBlocks2 := 4; hBlocks2 := 4  // reduction block size for peak pool; require power of 2
	LightState := false // off/on = false/true
	_ = LightState // --dev-- set to track gpio output state, and otherwise currently unused
	capture_diags := false  // --dev-- diagnostics mode; acquiare pico outputs from raspi

	// gpio config
	gpio10 := machine.GP10 // physical pin 14, physical pin 13 == gnd
	led := machine.LED
	gpio10.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// initialize iSpectRefReduced* and Create* loop memory
	fftPoints := buf_size/Tbins // e.g. for 1024 with 64 Tbins: (/ 1024 64) 16 points per fft; require power of 2
	if IsPow2(fftPoints) != true {
		panic("fft() requires power of 2 input size" + GetFunctionName(CreateU16SpectFromU16))
	}
	ref_init := make([]uint16, buf_size) // for allocation sizing only
	HammingFftPoints := Hamming(fftPoints) 
	U16SpectRef, _ := CreateU16SpectFromU16 ( ref_init, HammingFftPoints, Tbins, Fbins, buf_size, SpectThresh )
	iSpectRefReducedLight, iSpectRefReducedLight_PoolAvg := ReduceWordDetectCreateRef( U16SpectRef, Fbins, Tbins,
		vBlocks, hBlocks, vBlocks2, hBlocks2 )
	iSpectRefReducedDark, _  := ReduceWordDetectCreateRef( U16SpectRef, Fbins, Tbins,
		vBlocks, hBlocks, vBlocks2, hBlocks2 )	
	ref_init = nil
	
	// fmt.Printf("First 2 sounds set 'light' and 'dark' ref\n\r")
	loopCt := 0 
	for { // --ever--
		
		// --quiet-- fmt.Printf("Waiting for sound...") 
		// --quiet-- fmt.Printf("sound...") 
		uBuf := adc.Cap2Uint16(buf_size, sleep_time)
		if len(uBuf) < MinWordLen {
			flashOn(led); flashOn(led)
			continue
		}
		// --quiet-- fmt.Printf(" ct: %d\n\r", len(uBuf))
		// fmt.Println("--debug-- len(uBuf):", len(uBuf))
		// fmt.Printf("--debug-- uBuf:\n\r") // capture raw samples with minicom

		// process first two captures as ref words
		if loopCt < 2 {
			// verify time domain uBuf is not noise; 0xBFFF is 0.75 0xFFFF
			_, bIsNoise := NormalizeU16_ac_threshold(uBuf, 0xBFFF)
			if bIsNoise { // don't process and repeat this loop pass
				flashOn(led)
				continue
			}
			
			if loopCt == 0 {
				// create new light ref from initial capture
				U16SpectRef, bIsNoise = CreateU16SpectFromU16 ( uBuf, HammingFftPoints,
					Tbins, Fbins, buf_size, SpectThresh )
				iSpectRefReducedLight, iSpectRefReducedLight_PoolAvg =
					ReduceWordDetectCreateRef( U16SpectRef, Fbins, Tbins,
						vBlocks, hBlocks, vBlocks2, hBlocks2 )

				if capture_diags { // raspi diagnostics acquisition
					// create --uart out-- files for *_xt.dat, *_spect.dat, *_pool1/2.dat
					// '--' tagging embedded in uartHeader(); requires --eod-- to close file write	
					captureDiags( uBuf, U16SpectRef, iSpectRefReducedLight_PoolAvg, iSpectRefReducedLight)
				} // end if capture_diags 
			}
			if loopCt == 1 {
				// create new dark ref from initial capture
				U16SpectRef, bIsNoise = CreateU16SpectFromU16 ( uBuf, HammingFftPoints,
					Tbins, Fbins, buf_size, SpectThresh )
				iSpectRefReducedDark, _  = ReduceWordDetectCreateRef( U16SpectRef, Fbins, Tbins,
					vBlocks, hBlocks, vBlocks2, hBlocks2 )	
			}
		} // end if loopCt < 2
		loopCt++

		U16Spect, bIsNoise := CreateU16SpectFromU16 ( uBuf, HammingFftPoints,
			Tbins, Fbins, buf_size, SpectThresh )
		// fmt.Println("--debug-- U16Spect:", U16Spect[0][0:32],"\n\r")

		if bIsNoise {
			flashOn(led)
			continue
		}

		isLight := ReduceWordDetect(
			U16Spect, iSpectRefReducedLight, iSpectRefReducedDark, SpectThresh, buf_size, Fbins, Tbins,
			vBlocks, hBlocks, vBlocks2, hBlocks2 )

		// physical signifiers
		if loopCt < 3 {  // training
			if loopCt == 1 { 
				// first word; flash signifies training; note LoopCt was inc'd; trained 'light'
				// fmt.Println("loopCt == 1\n\r")
				flashOn(gpio10) // gpio10 flash and leave on
				LightState = true
			}
			if loopCt == 2 { // second word; flash to signifies training; trained 'dark'
				// fmt.Println("loopCt == 2\n\r")
				flashOff(gpio10) // gpio10 flash and leave off
				LightState = false
			}
		} else {  // not training
			if isLight == 1 {
				// fmt.Println("--d-- isLight:", isLight, "\n\r")
				gpio10.High()
				LightState = true
			}
			if isLight == 0 {
				gpio10.Low()
				LightState = false
			}
			// ReduceWordDetect() may return '3' or other to signify 'word not detected' 
		} // end if loopCt < 3
		
		// U16Spect = nil // --dev--
		// runtime.GC()   // --dev-- 
		
	} // end for --ever--
} // end main

// flashOn flashes the received gpio pin and leaves it in the on state
func flashOn( gpioPin machine.Pin ) {
	gpioPin.High()
	time.Sleep(time.Millisecond * 200) 
	gpioPin.Low()
	time.Sleep(time.Millisecond * 200) 
	gpioPin.High()
}

// flashOn flashes the received gpio pin and leaves it in the off state
func flashOff( gpioPin machine.Pin ) {
	gpioPin.Low()
	time.Sleep(time.Millisecond * 200) 
	gpioPin.High()
	time.Sleep(time.Millisecond * 200) 
	gpioPin.Low()
}

// uartHeader outputs Tag_file and 'filename' to stdout (uart)
// Transfer is ongoing until Tag_eot is sent to stdout (uart)
// Uart assumes receipt of Tag_eod to end the file start created here
// This code provides uart_xfr support on the raspi
func uartHeader( filename string) {
	fmt.Printf("........%s --%s--\n\r", Tag_file, filename)
	// transfer is ongoing
}

// uartFooter ends uart xfr started with uartHeader
func uartFooter() {
	fmt.Printf("%s\n\r", Tag_eot)
}

// captureDiags outputs the spectrogram and pooling arrays to stdout, intended
// for uart capture.  Arrays are wrapped with tags to assist parsing.
func captureDiags( uBuf []uint16, U16SpectRef [][]uint16,
	iSpectRefReducedLight_PoolAvg, iSpectRefReducedLight [][]int ) { 
	// create --uart out-- files for *_xt.dat, *_spect.dat, *_pool1/2.dat
	// '--' tagging embedded in uartHeader(); requires --eod-- to close file write
	uartHeader("file00_xt.dat") // u16 decimal 0-65535 (expt 2 16) 65536
	for _,v := range uBuf[:] {
		fmt.Println(v,"\n\r")
	}
	fmt.Println("--eod--","\n\r") // end of file00_xt.dat
	
	fmt.Println(Tag_file, "--file00_spect.dat--","\n\r")
	for i,_ := range U16SpectRef[:] {
		for j,_ := range U16SpectRef[0] {
			fmt.Printf("%d ", U16SpectRef[i][j])
		}
		fmt.Println("\n\r")
	}
	fmt.Println("--eod--","\n\r") // end of file00_spect.dat
	
	fmt.Println(Tag_file, "--file00_pool1.dat--","\n\r")
	for i,_ := range iSpectRefReducedLight_PoolAvg[:] {
		for j,_ := range iSpectRefReducedLight_PoolAvg[0] {
			fmt.Printf("%d ", iSpectRefReducedLight_PoolAvg[i][j])
		}
		fmt.Println("\n\r")
	}
	fmt.Println("--eod--","\n\r") // end of file00_pool1.dat
	
	fmt.Println(Tag_file, "--file00_pool2.dat--","\n\r")
	for i,_ := range iSpectRefReducedLight[:] {
		for j,_ := range iSpectRefReducedLight[0] {
			fmt.Printf("%d ", iSpectRefReducedLight[i][j])
		}
		fmt.Println("\n\r")
	}
	fmt.Println("--eod--","\n\r") // end of file00_pool2.dat
	uartFooter() // --eot--
} // end func captureDiags( uBuf []uint16, U16SpectRef [][]uint16,...

// Notes
// fmt.Println("--debug-- GetFunctionName() test :", GetFunctionName(U16HexList2GoIncludeVar))


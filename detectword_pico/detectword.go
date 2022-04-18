// @file TinyGo/sandbox_dw/detectword.go
// @date 2022.03.08
// @info detectword_pico specific functions

// Copyright 2022 RC Schuler. All rights reserved.
// Use of this source code is governed by a GNU V3
// license that can be found in the LICENSE file.

// @date 2022.03.15 added --warning-- CreateU16Spect clipping float %v to uint16 0
// @date 2022.03.15 copied and updated from sandbox_dw/detectword.go and commented --no pico-- 'exec' use
// @date 2022.03.16 added ReduceDetectWordCreateRef(); moved Ref processing out of loops
// @date 2022.03.25 changed ReduceWordDetect deltaLseDse from 200 to 0; better for live captured references;
//                  added 'deltaLseDseNoiseNeg/Pos' to detection scheme
// @date 2022.03.27 CreateU16SpectFromU16FromU16(); removed hex processing
// @date 2022.03.28 FftLogShift(); combines 20*math.Log10() and fft shift
// @date 2022.04.01 added NormalizeU16_ac_threshold() calls for sound level detection
// @date 2022.04.08 code cleanup; added bIsNoise as Create*FromU*() return
// @date 2022.04.09 tuning: deltaLseDseNoiseNeg/deltaLseDseNoisePos from -250/250 to -400/400
// @date 2022.04.13 commented all Print* for --prod-- mode; see --quiet--

// @build: tinygo flash -target=pico

package main

import (
	"fmt"
	"math"
	"os"
	// --raspi only-- "os/exec"
)

// --obs-- deprecated dev code for backards compatability; use for < v0.3 only 
// CreateSpectFileAndPlot outputs captured spectrogram to file and plots
// comment with // --no-pico-- for pico jobs
func CreateSpectFileAndPlot( varname string, Tsamp float64, Tbins, Fbins int, U16Spect [][]uint16 ) {
	captureFilename := fmt.Sprintf("file000_%s_spect.out",varname)     
	SpectrogramU16ToFile(captureFilename, Tbins, Fbins, U16Spect)      
	CreateOctaveSpect(captureFilename, Tsamp)                          
}

// --obs-- deprecated dev code for backards compatability; use for < v0.3 only 
func CreateOctaveSpect(captureFilename string, Tsamp float64) {
	// --raspi only-- osCmd := exec.Command("echo", "") // set osCmd and osErr types
	// --raspi only-- osErr := osCmd.Run()
	// --raspi only-- // fmt.Printf("--debug-- creating spectogram plot\n")
	// --raspi only-- osCmd = exec.Command("octave", "mfiles/plt_spect_nogui2.m", fmt.Sprintf("%f",Tsamp),
	// --raspi only-- fmt.Sprintf("%s", captureFilename))
	// --raspi only-- osErr = osCmd.Run()
	// --raspi only-- if osErr != nil {
	// --raspi only-- panic(fmt.Sprintf("CreateOcataveSpect exec error: %s",osErr))
	// --raspi only-- }
}

// --obs-- keeping this to document the various dev training words from mp3 files
// CreateGoIncludeVars reads in a hex representation of adc samples, assigns them to a const variable name
// and creates an include file. Mostly commented since Go const var include files only need updating 
// when reference word files are recaptured by adc
func CreateGoIncludeVars() {
	// create new include vars, e.g light_on, light_off
	// var sGoIncludeVar string // --dev-- only uncomment when new include files are created
	
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_light01", "light01")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_dark01", "dark01")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])

	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_rosie01", "rosie01")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])	
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_rosie02", "rosie02")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])	
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_angel01", "angel01")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])	
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_angel02", "angel02")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])	
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_light02", "light02")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])	
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_light03", "light03")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])	
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_light04", "light04")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])	
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_dark02", "dark02")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])	
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_dark03", "dark03")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])	
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_dark04", "dark04")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])

	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_mic_light01", "mic_light01")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_mic_light02", "mic_light02")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_mic_dark01", "mic_dark01")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])
	// sGoIncludeVar = U16HexList2GoIncludeVar("file00_mic_dark02", "mic_dark02")	
	// fmt.Println("--debug-- sGoIncludeVar:\n\r", sGoIncludeVar[0:64])		
	
	// display active varname, e.g. 'light_on', after U16HexList2GoIncludeVar() run;
	// These must be commented before go include file created and loaded in previous run
	// fmt.Println("--debug-- light_on: \n\r", light_on[0:64], "\n\r")
	// fmt.Println("--debug-- light_off: \n\r", light_off[0:64], "\n\r")
	// fmt.Println("light01: \n\r", light01[0:64], "\n\r")
	// fmt.Println("dark01: \n\r", dark01[0:64], "\n\r")
	// fmt.Println("rosie01: \n\r", rosie01[0:64], "\n\r")
	// fmt.Println("rosie02: \n\r", rosie02[0:64], "\n\r")
	// fmt.Println("angel01: \n\r", angel01[0:64], "\n\r")
	// fmt.Println("angel02: \n\r", angel02[0:64], "\n\r")
	// fmt.Println("light02: \n\r", light02[0:64], "\n\r")
	// fmt.Println("light03: \n\r", light03[0:64], "\n\r")
	// fmt.Println("light04: \n\r", light04[0:64], "\n\r")
	// fmt.Println("dark02: \n\r", dark02[0:64], "\n\r")
	// fmt.Println("dark03: \n\r", dark03[0:64], "\n\r")
	// fmt.Println("dark04: \n\r", dark04[0:64], "\n\r")
	
} // end func CreateGoIncludeVars() 

// FftLogShift consolidates fft shift and 20*math.Log10(); specific to CeaateU16SpectFromU16
func FftLogShift(fftReal []float64) (fftRealShift []float64) {
	fftRealShift = make([]float64, len(fftReal))
	mid := len(fftReal)/2
	for i,v := range fftReal[mid:] {
		fftRealShift[i] = 20.0*math.Log10(v)
	}
	for i:=0; i<mid; i++ {
		fftRealShift[mid+i]=20*math.Log10(fftReal[i])
	}
	return fftRealShift
} // end func FftLogShift(fftReal []float64) (fftRealShift []float64) {

// CreateU16SpectFromU16 receives 'u16Samples' time domain, and reusable buffer
// complexFloatArray, resizes tie domain to 'newsize',
// converts to 16 bit normalized 'ac' []int, loads into real part of []complex128 and
// create 'Tbins' fft's.  Number of frequency bins are resized to arbitray 'Fbins'.
// The collection of these fft's is stored as [][]u16Spect.  'newsize' must be a power of 2
// in place fft calculation.  Final log values below 'threshold' are set to zero on returned 'u16Spect'.
func CreateU16SpectFromU16 ( u16Samples []uint16, HammingFftPoints []float64, 
	Tbins, Fbins, newsize int, threshold uint16) (u16Spect [][]uint16, bIsNoise bool) {
	// create 'Tbins' ffts
	fftPoints := newsize/Tbins // e.g. for 2048: (/ 2048.0 64) 32.0 points per fft (require power of 2)
	if IsPow2(fftPoints) != true {
		panic("fft() requires power of 2 input size" + GetFunctionName(CreateU16SpectFromU16))
	}
	complexFloatArray := make( []complex128, fftPoints)
	u16Spect = make([][]uint16, Tbins) // second will be FbinFinal, allocated in main loop

	// noise filter threshold set to 0xBFFF which is 0.75 0xFFFF
	i16Samples, bIsNoise := NormalizeU16_ac_threshold(ResizeArrayUint16(u16Samples, newsize), 0xBFFF)
	if bIsNoise { // finish u16Spect allocation and return zeros
		for i,_ := range u16Spect {
			u16Spect[i] = make([]uint16, Fbins)
		}
		return u16Spect, bIsNoise // returning zeros indicating noise data set
	}
	// --obs-- u16Samples = nil
	// fmt.Println("--debug-- len normalized resized i16Samples:",len(i16Samples),"\n\r")
	// fmt.Println("--debug-- i16Samples:", i16Samples[0:8],"\n\r")
	// fmt.Println("--debug-- i16Samples:", i16Samples,"\n\r")

	lenComplexFloatArray := len(complexFloatArray)
	lenI16Samples := len(i16Samples)
	for i:=0; i< Tbins; i++ {
		for j:=0; j<lenComplexFloatArray; j++ {
			if i*fftPoints+j<lenI16Samples {
				complexFloatArray[j] = complex(
					HammingFftPoints[j] * float64(i16Samples[i*fftPoints+j]), // Hamming * Sample
					// float64(i16Samples[i*fftPoints+j]), // Sample
					0.0)
			} else {
				complexFloatArray[j] = complex(0.0,0.0) // zero pad
			}
		} // end for j:=0; j<lenComplexFloatArray; j++ 

		err := FFT(complexFloatArray)
		if err != nil {
			panic(err)
		}
		// if i<2 { fmt.Println("--debug-- CreateU16Spect() FFT():",complexFloatArray[0:4],"\n\r") }
		fftReal := Magnitude(complexFloatArray)
		fftRealShift := FftLogShift(fftReal) // --dev-- 20*math.Log10() and fft shift
		// --obs-- fftReal = nil
		
		// allocate Fbin dimension of u16Spect, apply threshold, and load return values
		u16Loader := make([]uint16, fftPoints)
		for k,v := range fftRealShift {
			if v < 0.0 {
				fmt.Fprintf(os.Stderr, "--warning-- CreateU16Spect clipped float %v to uint16 0\n\r", v)
				v = 0.0
			}
			if uint16(int(v)) < threshold {
				u16Loader[k] = threshold
				
			} else {
				u16Loader[k] = uint16(int(v))

			}
		}
		// --obs-- fftRealShift = nil
		// fmt.Println("--debug-- u16Loader:", u16Loader)
		u16Spect[i] = ResizeArrayUint16(u16Loader, Fbins)
		// --obs-- u16Loader = nil
		// fmt.Println("--debug-- u16Spect[i]:", u16Spect[i])

	} // end for i:=0; i< Tbins; i++

	return u16Spect, bIsNoise
} // end func CreateU16SpectFromU16

// --obs-- deprecated dev code for backards compatability; use for < v0.3 only 
func CreateU16SpectFromU16_sync ( u16Samples []uint16, complexFloatArray []complex128, HammingFftPoints []float64, 
	Tbins, Fbins, newsize int, threshold uint16) (u16Spect [][]uint16) {
	// create 'Tbins' ffts
	maxCaptureSize := newsize  // e.g. for 8192: (/ 8192.0 128) 64.0 points per fft (require power of 2)
	fftPoints := maxCaptureSize/Tbins
	if IsPow2(fftPoints) != true {
		panic("fft() requires power of 2 input size" + GetFunctionName(CreateU16SpectFromU16))
	}
	// --obs-- complexFloatArray := make( []complex128, fftPoints)
	u16Spect = make([][]uint16, Tbins) // second will be FbinFinal, allocated in main loop
	// --obs-- u16Samples := StringHexBytes2Uint16( sU16Hex )
	// --obs-- sU16Hex=""
	// fmt.Println("--debug-- len u16Samples:",len(u16Samples), "\n\r")

	// noise filter threshold set to 0xBFFF which is 0.75 0xFFFF
	i16Samples, bIsNoise := NormalizeU16_ac_threshold(ResizeArrayUint16(u16Samples, newsize), 0xBFFF)
	if bIsNoise { // finish u16Spect allocation and return zeros
		for i,_ := range u16Spect {
			u16Spect[i] = make([]uint16, Fbins)
		}
		return u16Spect // returning zeros indicating noise data set
	}
	u16Samples=nil
	// fmt.Println("--debug-- len normalized resized i16Samples:",len(i16Samples),"\n\r")
	// fmt.Println("--debug-- i16Samples:", i16Samples[0:8],"\n\r")
	// fmt.Println("--debug-- i16Samples:", i16Samples,"\n\r")

	lenComplexFloatArray := len(complexFloatArray)
	lenI16Samples := len(i16Samples)
	for i:=0; i< Tbins; i++ {
		fftReal := HammingFftPoints // fftReal multi-tasks between fft result and Hamming multiplier
		for j:=0; j<lenComplexFloatArray; j++ {
			if i*fftPoints+j<lenI16Samples {
				complexFloatArray[j] = complex(
					fftReal[j] * float64(i16Samples[i*fftPoints+j]), // Hamming * Sample
					// float64(i16Samples[i*fftPoints+j]), // Sample
					0.0)
			} else {
				complexFloatArray[j] = complex(0.0,0.0) // zero pad
			}
		} // end for j:=0; j<lenComplexFloatArray; j++ 

		// if i<2 { fmt.Println("--debug-- #####\n\r") }
		// if i<2 { fmt.Println("--debug-- CreateU16Spect() complexFloatArray[0:4]:", complexFloatArray[0:4],"\n\r") }
		err := FFT(complexFloatArray)
		if err != nil {
			panic(err)
		}
		// if i<2 { fmt.Println("--debug-- CreateU16Spect() FFT():",complexFloatArray[0:4],"\n\r") }
		fftReal = Magnitude(complexFloatArray)

		fftRealShift := FftLogShift(fftReal) // --dev-- 20*math.Log10() and fft shift
		fftReal=nil
		
		// allocate Fbin dimension of u16Spect, apply threshold, and load return values
		u16Loader := make([]uint16, fftPoints)
		for k,v := range fftRealShift {
			if v < 0.0 {
				fmt.Fprintf(os.Stderr, "--warning-- CreateU16Spect clipped float %v to uint16 0\n\r", v)
				v = 0.0
			}
			if uint16(int(v)) < threshold {
				u16Loader[k] = threshold
				
			} else {
				u16Loader[k] = uint16(int(v))

			}
		}
		fftRealShift=nil
		// fmt.Println("--debug-- u16Loader:", u16Loader)
		u16Spect[i] = ResizeArrayUint16(u16Loader, Fbins)		
		// fmt.Println("--debug-- u16Spect[i]:", u16Spect[i])

	} // end for i:=0; i< Tbins; i++

	return u16Spect
} // end func CreateU16SpectFromU16_sync

// --obs-- deprecated dev code for backards compatability; use for < v0.3 only 
// SpectrogramU16ToFile receives outnname for file, Tbins, Fbins, and [][]uint16 Spect, and
// writes data to file 'outname' ins space separated mxn octave format.
// Runs on raspi; default pico has no file system.
func SpectrogramU16ToFile(outname string, Tbins, Fbins int, Spect [][]uint16) {
	fileOut, err := os.Create(outname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic: %s openning %s\n",
			GetFunctionName(SpectrogramU16ToFile), outname)
		panic(err)
	}
	defer fileOut.Close()
	for i:=0; i<Tbins; i++ {
		for j:=0; j<Fbins; j++ {
			fmt.Fprintf(fileOut,"%v ", Spect[i][j])
		}
		fmt.Fprintf(fileOut,"\n")
	}
} // end func SpectrogramU16ToFile(outname string, Tbins, Fbins int, Spect [][]uint16) 

// ReduceWordDetect resolves U16SpectRef and U16Spect into word 'Light' or 'Dark'
// U16SpectRef have been reduced to 'iSpectReducedLight/Dark' before call
// --prod-- tuned with SpectThresh=50, vBlocks=8, hBlocksk=8 (avg), vBlocks2=4, hBlocks2=4 (peak),
// Fbins=64, Tbins=64, buf_size=1024, Tsamp=166us; deltaLseDse=0
func ReduceWordDetect(
	U16Spect [][]uint16, iSpectRefReducedLight, iSpectRefReducedDark [][]int, SpectThresh uint16, buf_size, Fbins, Tbins,
	vBlocks, hBlocks, vBlocks2, hBlocks2 int ) (isLight int ) {
	
	deltaLseDse := 0 // detla (lse-dse) decision point; was < 200 == 'Light'
	deltaLseDseNoiseNeg := -400 // 20220409 was -250/250
	deltaLseDseNoisePos :=  400
	
	rows0 := Fbins; cols0 := Tbins    // org array size // --dev-- --was-- Fbins/2
	rows1 := rows0/vBlocks; cols1 := cols0/hBlocks // working subslice size
	// fmt.Printf("--d-- rows0: %d, cols0: %d\n\r", rows0, cols0)
	// fmt.Printf("--d-- rows1: %d, cols1: %d\n\r", rows1, cols1)

	// detect light 
	// first reduction by average value

	// --dev-- for Fbins/2 iSpectReduced :=  ReduceUint16ToIntArrayAvg( U16Spect[rows0:][:], rows1, cols1 )
	iSpectReduced :=  ReduceUint16ToIntArrayAvg( U16Spect[:][:], rows1, cols1 )
	// --dev-- iSpectReduced :=  ReduceUint16ToIntArrayPeak( U16Spect[rows0:][:], rows1, cols1 )
	// fmt.Println("--d-- iSpectReduced   ", iSpectReduced )
	
	// second reduction by peak value
	rows0 = len(iSpectReduced); cols0 = len(iSpectReduced[0])    // org array size
	rows1 = rows0/vBlocks2; cols1 = cols0/hBlocks2 // working subslice size
 	// fmt.Printf("--d-- rows0: %d, cols0: %d\n\r", rows0, cols0)
	// fmt.Printf("--d-- rows1: %d, cols1: %d\n\r", rows1, cols1)		
	iSpectReduced = ReduceIntArrayPeak( iSpectReduced, rows1, cols1 )
	// fmt.Println("--d-- iSpectReduced   ", iSpectReduced )
	
	// calc and print square err vs 'light'
	for i,_ := range iSpectRefReducedLight {
		for j,_ := range iSpectRefReducedLight[0] {
			iSpectReduced[i][j] =
				int(math.Floor(
					math.Pow(float64(iSpectRefReducedLight[i][j])-float64(iSpectReduced[i][j]), 2)))
		}
	}
	lse := SliceSumInt(iSpectReduced) 
	// --quiet-- fmt.Println("LSE:", iSpectReduced, lse, "\n\r ") // light sq err
	
	// detect dark
	// avg reduction word detection
	rows0 = Fbins; cols0 = Tbins    // org array size // --dev-- --was-- Fbins/2
	rows1 = rows0/vBlocks; cols1 = cols0/hBlocks // working subslice size
	// fmt.Printf("--d-- rows0: %d, cols0: %d\n\r", rows0, cols0)
	// fmt.Printf("--d-- rows1: %d, cols1: %d\n\r", rows1, cols1)		

	// first reduce by average value
	// --dev-- for Fbins/2: iSpectReduced =  ReduceUint16ToIntArrayAvg( U16Spect[rows0:][:], rows1, cols1 )
	iSpectReduced =  ReduceUint16ToIntArrayAvg( U16Spect[:][:], rows1, cols1 )
	// --dev-- iSpectReduced =  ReduceUint16ToIntArrayPeak( U16Spect[rows0:][:], rows1, cols1 )
	// fmt.Println("--d-- iSpectReduced   ", iSpectReduced )
	
	// reduce a second time by peak value
	rows0 = len(iSpectReduced); cols0 = len(iSpectReduced[0])    // org array size
	rows1 = rows0/vBlocks2; cols1 = cols0/hBlocks2 // working subslice size
	// fmt.Printf("--d-- rows0: %d, cols0: %d\n\r", rows0, cols0)
	// fmt.Printf("--d-- rows1: %d, cols1: %d\n\r", rows1, cols1)		
	// --dev-- iSpectReduced = ReduceIntArrayAvg( iSpectReduced, rows1, cols1 )
	iSpectReduced = ReduceIntArrayPeak( iSpectReduced, rows1, cols1 )
	// fmt.Println("--d-- iSpectReduced   ", iSpectReduced )
	
	// calc and print square err vs 'dark'
	for i,_ := range iSpectRefReducedDark {
		for j,_ := range iSpectRefReducedDark[0] {
			iSpectReduced[i][j] =
				int(math.Floor(
					math.Pow(float64(iSpectRefReducedDark[i][j])-float64(iSpectReduced[i][j]), 2)))
		}
	}
	dse := SliceSumInt(iSpectReduced) // dark sq err
	// --quiet-- fmt.Println("DSE:", iSpectReduced, dse, "del", lse-dse, "\n\r" ) // dark sq err
	// decision:
	lseMinusDse := lse-dse
	isLight = 3 // set to 'noise detected'
	if (lseMinusDse <= deltaLseDse) && (lseMinusDse > deltaLseDseNoiseNeg) { 
		// --quiet-- fmt.Println("\"Light\"", "\n\r")
		isLight = 1
	}
	if (lseMinusDse > deltaLseDse) && (lseMinusDse < deltaLseDseNoisePos) { 	
		// --quiet-- fmt.Println("\"Dark\"", "\n\r")
		isLight = 0
	}

	return isLight
} // end ReduceWordDetect

// ReduceWordDetectCreateRef provides a separate reduction function for reference words and
// returns both the final reduction, and the intermediate pool1 state for diagnostics
func ReduceWordDetectCreateRef( U16SpectRef [][]uint16, Fbins, Tbins int,
	vBlocks, hBlocks, vBlocks2, hBlocks2 int ) (i16SpectRefReduced, i16SpectRefReducedPoolAvg [][]int  ) {
	rows0 := Fbins; cols0 := Tbins    // org array size // --dev-- --was-- Fbins/2
	rows1 := rows0/vBlocks; cols1 := cols0/hBlocks // working subslice size

	// first reduction by average value
	// --dev-- for Fbins/2: iSpectRefReduced := ReduceUint16ToIntArrayAvg( U16SpectRef[rows0:][:], rows1, cols1 )
	iSpectRefReduced := ReduceUint16ToIntArrayAvg( U16SpectRef[:][:], rows1, cols1 )
	// --dev-- iSpectRefReduced := ReduceUint16ToIntArrayPeak( U16SpectRef[rows0:][:], rows1, cols1 )
	// fmt.Println("--d-- iSpectRefReduced", iSpectRefReduced )

	// create intermediate copy for return diags
	lenIspectRef := len(iSpectRefReduced)
	lenIspectRef0 := len(iSpectRefReduced[0])
	iSpectRefReducedPoolAvg := make([][]int, lenIspectRef)
	for i,_ := range iSpectRefReducedPoolAvg { // allocate new memory
		iSpectRefReducedPoolAvg[i] = make ([]int, len(iSpectRefReduced[0]))
	}
	for i:=0; i<lenIspectRef; i++ { // perform copy 
		for j:=0; j<lenIspectRef0; j++ {
			iSpectRefReducedPoolAvg[i][j] = iSpectRefReduced[i][j]
		}
	}

	// second reduction by peak value
	rows0 = len(iSpectRefReduced); cols0 = len(iSpectRefReduced[0])    // org array size
	rows1 = rows0/vBlocks2; cols1 = cols0/hBlocks2 // working subslice size
 	// fmt.Printf("--d-- rows0: %d, cols0: %d\n\r", rows0, cols0)
	// fmt.Printf("--d-- rows1: %d, cols1: %d\n\r", rows1, cols1)		
	// --dev-- iSpectRefReduced = ReduceIntArrayAvg( iSpectRefReduced, rows1, cols1 )
	iSpectRefReduced = ReduceIntArrayPeak( iSpectRefReduced, rows1, cols1 )
	// fmt.Println("--d-- iSpectRefReduced", iSpectRefReduced )

	return iSpectRefReduced, iSpectRefReducedPoolAvg
} // end ReduceWordDetectCreateRef

// @file TinyGo/sandbox_dw/utils.go
// @date 2022.03.08
// @info detectword utils

// @info contains foss go-fft code tagged with https://github.com/ledyba/go-fft/blob/master/LICENSE

// Copyright 2022 RC Schuler. All rights reserved.
// Use of this source code is governed by a GNU V3
// license that can be found in the LICENSE file.

// @date 2022.03.14 additions from reduce_array_avg dev
// @date 2022.04.01 added Normalize_ac_threshold()

// @build: include file
package main

import (
	"fmt"
	"math"
	"math/bits"
	"os"
	"bufio"
	"runtime"
	"reflect"
	"strconv"
)
// BetweenFloat returns true if f is between min and max, else false
func BetweenFloat(f, min, max float64) bool {
         if (f >= min) && (f <= max) {
                 return true
         } else {
                 return false
         }
 }

// CreateHamming returns a Hamming (cos) window of size n.
func Hamming(n int) []float64 {
	w := make([]float64, n)

	if n == 1 {
		w[0] = 1
	} else {
		N := n - 1
		weight := math.Pi * 2 / float64(N)
		for n := 0; n <= N; n++ {
			w[n] = 0.54 - 0.46*math.Cos(weight*float64(n))
		}
	}
	return w
}

// ReverseU16Array receive a []Uint16 and return it byte reversed
func ReverseU16rray(arr []uint16) []uint16{
   for i, j := 0, len(arr)-1; i<j; i, j = i+1, j-1 {
      arr[i], arr[j] = arr[j], arr[i]
   }
   return arr
}

// NormalizeU16_ac() Normalize a slice 0000->FFFF, then remove avg (dc), return as []int 
func NormalizeU16_ac(data []uint16) []int {
	// use float operations
	fdata := make([]float64, len(data))
	for i,v := range data {
		fdata[i]=float64(v)
	}
	// subtract min value, divide by max value after shift of each element
	mi :=fdata[0] // min accum
	mx :=fdata[0] // max accum
	for _, e := range fdata {
		if e < mi {
			mi = e 
		}
		if e > mx  {
			mx = e 
		} 
	}
	avg := 0.0 
	for i, e := range fdata {
		fdata[i]=(e-mi)/(mx-mi) * float64(0xffff)
		avg += fdata[i]
 	}
	avg = avg / float64(len(fdata))
	idata := make([]int, len(fdata))
	for i, e := range fdata {
		idata[i] = int(e-avg)
	}
	return idata
}

// NormalizeU16_ac_threshold() Normalize a slice 0000->FFFF, then remove avg (dc), return as []int,
// reject data set which never exceeds 'dataThreshold', returning with 'bIsNoise' set true;
// also reject as noise high end clipping (>0xFFFD)
func NormalizeU16_ac_threshold(data []uint16, dataThreshold uint16) (idata[] int, bIsNoise bool) {
	// use float operations
	fdata := make([]float64, len(data))
	for i,v := range data {
		fdata[i]=float64(v)
	}
	idata = make([]int, len(fdata))

	// subtract min value, divide by max value after shift of each element
	mi :=fdata[0] // min accum
	mx :=fdata[0] // max accum
	for _, e := range fdata {
		if e < mi {
			mi = e 
		}
		if e > mx  {
			mx = e 
		} 
	}

	if uint16(mx) < dataThreshold || uint16(mx) > 0xFFF0 {  // noise data set, return previously allocated zeros
		// fmt.Println("--debug--", "mx:", mx, "\n\r")
		return idata, true // bIsNoise is true
 	}
		
	avg := 0.0 
	for i, e := range fdata {
		fdata[i]=(e-mi)/(mx-mi) * float64(0xffff)
		avg += fdata[i]
 	}
	avg = avg / float64(len(fdata))
	for i, e := range fdata {
		idata[i] = int(e-avg)
	}
	return idata, false // bIsNoise is false
}

// NormalizeU16() Normalize a slice 0000->FFFF, return as []uint16
// Changes from NormalizeU16_ac() tagged with --noac--
func NormalizeU16(data []uint16) []uint16 { // --noac--
	// use float operations
	fdata := make([]float64, len(data))
	for i,v := range data {
		fdata[i]=float64(v)
	}
	// subtract min value, divide by max value after shift of each element
	mi :=fdata[0] // min accum
	mx :=fdata[0] // max accum
	for _, e := range fdata {
		if e < mi {
			mi = e 
		}
		if e > mx  {
			mx = e 
		} 
	}
	avg := 0.0  
	for i, e := range fdata {
		fdata[i]=(e-mi)/(mx-mi) * float64(0xffff)
		avg += fdata[i]
 	}
	avg = avg / float64(len(fdata))
	udata := make([]uint16, len(fdata)) // --noac--
	for i, e := range fdata { 
		udata[i] = uint16(e-avg)    // --noac--
	}                      
	return udata
}

// U16HexList2GoIncludeVar receives a 'filename' and an include 'varname';
// creates the output file 'include_<varname>.go' as a side effect.
// This function is for creating diagnostics data sets with vars from
// generated include files.
// (note: not "\n\r" newlines; this runs on rapi host and not pico)
// example output:
// @file Tinygo/detectword_pico/
// @date 2022.02.28

// package main
// // --start--
// const ref_light_on = "84e0706071..."
// // --stop--
func U16HexList2GoIncludeVar( filename, varname string) (GoInclude string) {
	file, err := os.Open(filename) 
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic: %s openning %s\n",
			GetFunctionName(U16HexList2GoIncludeVar), filename)
		panic(err)
	}
	// initial header
	// --timestamp-- disabled by quoting 'time.Now()' sprintf; replace <quote> with " for timestamps;
	// --timestamp-- preventing updates to include files in --dev-- mode
	GoInclude = "\n// @file include_" + varname + ".go\n" + "// @date " + "fmt.Sprintf(<quote>%s<quote>,time.Now())[:16]" + "\n\n" + "package main\n" + "// --start--\n" + "const " + varname + " = \""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		GoInclude = GoInclude + line
	}
	// add footer
	GoInclude = GoInclude + "\"\n// --stop--\n"

	// write output include*.go file
	outname := fmt.Sprintf("include_%s.go", varname)
	fmt.Println("outfile:", outname, "\n")
	fileOut, err := os.Create(outname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic: %s openning %s\n",
			GetFunctionName(U16HexList2GoIncludeVar), filename)
		panic(err)
	}
	defer fileOut.Close()
	fmt.Fprintf(fileOut,"%s", GoInclude)
	return GoInclude
} // end func U16HexList2GoIncludeVar( filename, varname string) (GoInclude string) 

// U16HexList2String is U16HexLIst2GoIncludeVar without headers, return string stream hex data only
func U16HexList2String( filename, varname string) (GoInclude string) {
	file, err := os.Open(filename) 
	if err != nil {
		fmt.Fprintf(os.Stderr, "panic: %s openning %s\n",
			GetFunctionName(U16HexList2GoIncludeVar), filename)
		panic(err)
	}
	// initial header
	// --timestamp-- disabled by quoting 'time.Now()' sprintf; replace <quote> with " for timestamps;
	// --timestamp-- preventing updates to include files in --dev-- mode
	// --no-header-- GoInclude = "\n// @file include_" + varname + ".go\n" + "// @date " + "fmt.Sprintf(<quote>%s<quote>,time.Now())[:16]" + "\n\n" + "package main\n" + "// --start--\n" + "const " + varname + " = \""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		GoInclude = GoInclude + line
	}
	// add footer
	// --no-header-- GoInclude = GoInclude + "\"\n// --stop--\n"

	// --no-header-- // write output include*.go file
	// --no-header-- outname := fmt.Sprintf("include_%s.go", varname)
	// --no-header-- fmt.Println("--debug-- outfile:", outname, "\n")
	// --no-header-- fileOut, err := os.Create(outname)
	// --no-header-- if err != nil {
	// --no-header-- fmt.Fprintf(os.Stderr, "panic: %s openning %s\n",
	// --no-header-- GetFunctionName(U16HexList2GoIncludeVar), filename)
	// --no-header-- panic(err)
	// --no-header-- }
	// --no-header-- defer fileOut.Close()
	// --no-header-- fmt.Fprintf(fileOut,"%s", GoInclude)
	return GoInclude
} // end func U16HexList2String( filename, varname string) (GoInclude string) 

// GetFunctionName of passed function 'name' string
func GetFunctionName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// StringHexBytes2Uint16 takes a string of hex bytes representing a stream of
// uint16 values (e.g. "8500ffff..."). The array of uint16 values is returned
// as u16Bytes.  Passing a stream length not divisible by 4 is a fatal error.
// @todo place in pkg; dsp? From learn/hex_bytes_to_uint16 
func StringHexBytes2Uint16( sHexBytes string ) ( u16Bytes []uint16) {
	const step = 4
	if len(sHexBytes)%step != 0 {
		panic(fmt.Sprintf("Bad hex stream length (%d)\n\rFirst4: %s\n\rLast4: %s",
			len(sHexBytes),sHexBytes[0:4],sHexBytes[len(sHexBytes)-4:]))
	}
	numU16Bytes := len(sHexBytes)/4
	// fmt.Println("--debug-- numU16Bytes", numU16Bytes, "\n\r" )
	u16Bytes = make([]uint16, numU16Bytes)
	byteCount := 0
	for i:=0; i<len(sHexBytes)-3; i=i+step {
		// fmt.Println("--debug-- sHexBytes[i:i+4]:", sHexBytes[i:i+4], "\n\r")
		iValue, errParse := strconv.ParseInt(fmt.Sprintf("0x%s",sHexBytes[i:i+4]), 0, 64) // hex string to int
		if errParse != nil {
			panic(fmt.Sprintf("Hex string to int conversion error code %v",errParse))
		}
		u16Value := uint16(iValue) // int to uint16
		// fmt.Println("--debug-- u16Value:", u16Value, "\n\r")
		u16Bytes[byteCount] = u16Value
		byteCount = byteCount + 1
	} // end for i:=0; i<len(sHexBytes)-3; i=i+step 

	
	return u16Bytes
} // end func StringHexBytes2Uint16( sHexBytes ) ( u16Bytes []uint16) 

// Magnitude() returns the float64 Magnitude of a complex128 arg
func Magnitude(data []complex128) []float64 {
    magVals := make([]float64, len(data))
    for i, comp := range data {
        rel := math.Pow(real(comp), 2)
        img := math.Pow(imag(comp), 2)
        magVals[i] = math.Sqrt(rel + img)
    }
    return magVals
}

// ResizeArrayUnit16 receives a uint16 array and resizes
// it to new len 'n'. New elements based on rounded percentage 
// position within file, no interpolation. Returns the new
// uint16 array
func ResizeArrayUint16(u0 []uint16, n int) []uint16 {
	n0 := len(u0)
	u1 := make([]uint16, n)
	for i,_ := range u1 {
		pct := float64(i)/float64(n)
		u0Pos := int(math.Floor(pct*float64(n0)))
		u1[i] = u0[u0Pos]
	}
	return u1
} // end ResizeArrayUint16

// IsPow2 returns true if N is a perfect power of 2 (1, 2, 4, 8, ...) and false otherwise.
// Algorithm from: https://graphics.stanford.edu/~seander/bithacks.html#DetermineIfPowerOf2
// https://github.com/ledyba/go-fft/blob/master/LICENSE
func IsPow2(N int) bool {
	if N == 0 {
		return false
	}
	return (uint64(N) & uint64(N-1)) == 0
}

// NextPow2 returns the smallest power of 2 >= N.
// https://github.com/ledyba/go-fft/blob/master/LICENSE
func NextPow2(N int) int {
	if N == 0 {
		return 1
	}
	return 1 << uint64(bits.Len64(uint64(N-1)))
}

// ZeroPad pads x with 0s at the end into a new array of length N.
// This does not alter x, and creates an entirely new array.
// This should only be used as a convience function, and isn't meant for performance.
// You should call this as few times as possible since it does potentially large allocations.
// https://github.com/ledyba/go-fft/blob/master/LICENSE
func ZeroPad(x []complex128, N int) []complex128 {
	y := make([]complex128, N)
	copy(y, x)
	return y
}

// ZeroPadToNextPow2 pads x with 0s at the end into a new array of length 2^N >= len(x)
// This does not alter x, and creates an entirely new array.
// This should only be used as a convience function, and isn't meant for performance.
// You should call this as few times as possible since it does potentially large allocations.
// https://github.com/ledyba/go-fft/blob/master/LICENSE
func ZeroPadToNextPow2(x []complex128) []complex128 {
	N := NextPow2(len(x))
	y := make([]complex128, N)
	copy(y, x)
	return y
}

// Float64ToComplex128Array converts a float64 array to the equivalent complex128 array
// using an imaginary part of 0.
// https://github.com/ledyba/go-fft/blob/master/LICENSE
func Float64ToComplex128Array(x []float64) []complex128 {
	y := make([]complex128, len(x))
	for i, v := range x {
		y[i] = complex(v, 0)
	}
	return y
}

// Complex128ToFloat64Array converts a complex128 array to the equivalent float64 array
// taking only the real part.
func Complex128ToFloat64Array(x []complex128) []float64 {
	y := make([]float64, len(x))
	for i, v := range x {
		y[i] = real(v)
	}
	return y
}

// RoundFloat64Array calls math.Round on each entry in x, changing the array in-place
// https://github.com/ledyba/go-fft/blob/master/LICENSE
func RoundFloat64Array(x []float64) {
	for i, v := range x {
		x[i] = math.Round(v)
	}
}

// ReduceUint16ArrayAvg reduces arr0 using 'vSliceSize' x 'hSliceSizse' steps
// over arr0; elements of arr1 contain the avg value for a corresponding
// arr0 block; requires power of 2 slice params for valid size reduction
func ReduceUint16ArrayAvg(arr0 [][]uint16, vSliceSize, hSliceSize int) (arr1 [][]uint16) {
	lenArr0 := len(arr0)
	lenArr00 := len(arr0[0])
	fArr0 := make([][]float64, lenArr0)
	for i,_ := range arr0 {
		fArr0[i] = make([]float64, lenArr00)
	}
	for i,_ := range arr0 {
		for j,_ := range arr0[0] {
			fArr0[i][j] = float64(arr0[i][j])
		}
	}
	vReduced := lenArr0/vSliceSize
	hReduced := lenArr00/hSliceSize
	arr1 = make([][]uint16, vReduced)
	for i,_ := range arr1 {
		arr1[i] = make([]uint16, hReduced)
	}
	fArr1 := ReduceFloat64ArrayAvg(fArr0, vSliceSize, hSliceSize)
	// fmt.Println("--d-- len fArr1", len(fArr1))
	// fmt.Println("--d-- len fArr1[0]", len(fArr1[0]))
	for i,_ := range fArr1 {
		for j,_ := range fArr1[0] {
			arr1[i][j] = uint16(fArr1[i][j])
		}
	}
	return arr1
} // end func ReduceUint16ArrayAvg

// ReduceIntUint16ToIntArrayAvg reduces arr0 using 'vSliceSize' x 'hSliceSizse' steps
// over arr0; elements of arr1 contain the avg value for a corresponding
// arr0 block; requires power of 2 slice for valid size reduction;
// receives uint16 array and returns int array
func ReduceUint16ToIntArrayAvg(arr0 [][]uint16, vSliceSize, hSliceSize int) (arr1 [][]int) {
	lenArr0 := len(arr0)
	lenArr00 := len(arr0[0])
	fArr0 := make([][]float64, lenArr0)
	for i,_ := range arr0 {
		fArr0[i] = make([]float64, lenArr00)
	}
	for i,_ := range arr0 {
		for j,_ := range arr0[0] {
			fArr0[i][j] = float64(arr0[i][j])
		}
	}
	vReduced := lenArr0/vSliceSize
	hReduced := lenArr00/hSliceSize
	arr1 = make([][]int, vReduced)
	for i,_ := range arr1 {
		arr1[i] = make([]int, hReduced)
	}
	fArr1 := ReduceFloat64ArrayAvg(fArr0, vSliceSize, hSliceSize)
	// fmt.Println("--d-- len fArr1", len(fArr1))
	// fmt.Println("--d-- len fArr1[0]", len(fArr1[0]))
	for i,_ := range fArr1 {
		for j,_ := range fArr1[0] {
			arr1[i][j] = int(fArr1[i][j])
		}
	}
	return arr1
} // end func ReduceUint16ToIntArrayAvg

// ReduceIntUint16ToIntArrayPeak reduces arr0 using 'vSliceSize' x 'hSliceSizse' steps
// over arr0; elements of arr1 contain the peak value for a corresponding
// arr0 block; requires power of 2 slice params for valid size reduction;
// receives a uint16 array and returns an int array
func ReduceUint16ToIntArrayPeak(arr0 [][]uint16, vSliceSize, hSliceSize int) (arr1 [][]int) {
	lenArr0 := len(arr0)
	lenArr00 := len(arr0[0])
	fArr0 := make([][]float64, lenArr0)
	for i,_ := range arr0 {
		fArr0[i] = make([]float64, lenArr00)
	}
	for i,_ := range arr0 {
		for j,_ := range arr0[0] {
			fArr0[i][j] = float64(arr0[i][j])
		}
	}
	vReduced := lenArr0/vSliceSize
	hReduced := lenArr00/hSliceSize
	arr1 = make([][]int, vReduced)
	for i,_ := range arr1 {
		arr1[i] = make([]int, hReduced)
	}
	fArr1 := ReduceFloat64ArrayPeak(fArr0, vSliceSize, hSliceSize)
	// fmt.Println("--d-- len fArr1", len(fArr1))
	// fmt.Println("--d-- len fArr1[0]", len(fArr1[0]))
	for i,_ := range fArr1 {
		for j,_ := range fArr1[0] {
			arr1[i][j] = int(fArr1[i][j])
		}
	}
	return arr1
} // end func ReduceUint16ToIntArrayPeak

// ReduceIntArrayAvg reduces arr0 using 'vSliceSize' x 'hSliceSizse' steps
// over arr0; elements of arr1 contain the avg value for a corresponding
// arr0 block; requires power of 2 slice params for valid size reduction;
func ReduceIntArrayAvg(arr0 [][]int, vSliceSize, hSliceSize int) (arr1 [][]int) {
	lenArr0 := len(arr0)
	lenArr00 := len(arr0[0])
	fArr0 := make([][]float64, lenArr0)
	for i,_ := range arr0 {
		fArr0[i] = make([]float64, lenArr00)
	}
	for i,_ := range arr0 {
		for j,_ := range arr0[0] {
			fArr0[i][j] = float64(arr0[i][j])
		}
	}
	vReduced := lenArr0/vSliceSize
	hReduced := lenArr00/hSliceSize
	arr1 = make([][]int, vReduced)
	for i,_ := range arr1 {
		arr1[i] = make([]int, hReduced)
	}
	fArr1 := ReduceFloat64ArrayAvg(fArr0, vSliceSize, hSliceSize)
	// fmt.Println("--d-- len fArr1", len(fArr1))
	// fmt.Println("--d-- len fArr1[0]", len(fArr1[0]))
	for i,_ := range fArr1 {
		for j,_ := range fArr1[0] {
			arr1[i][j] = int(fArr1[i][j])
		}
	}
	return arr1
} // end func ReduceIntArrayAvg

// ReduceIntArrayPeak reduces arr0 using 'vSliceSize' x 'hSliceSizse' steps
// over arr0; elements of arr1 contain the peak value for a corresponding
// arr0 block; requires power of 2 slice params for valid size reduction
func ReduceIntArrayPeak(arr0 [][]int, vSliceSize, hSliceSize int) (arr1 [][]int) {
	lenArr0 := len(arr0)
	lenArr00 := len(arr0[0])
	fArr0 := make([][]float64, lenArr0)
	for i,_ := range arr0 {
		fArr0[i] = make([]float64, lenArr00)
	}
	for i,_ := range arr0 {
		for j,_ := range arr0[0] {
			fArr0[i][j] = float64(arr0[i][j])
		}
	}
	vReduced := lenArr0/vSliceSize
	hReduced := lenArr00/hSliceSize
	arr1 = make([][]int, vReduced)
	for i,_ := range arr1 {
		arr1[i] = make([]int, hReduced)
	}
	fArr1 := ReduceFloat64ArrayPeak(fArr0, vSliceSize, hSliceSize)
	// fmt.Println("--d-- len fArr1", len(fArr1))
	// fmt.Println("--d-- len fArr1[0]", len(fArr1[0]))
	for i,_ := range fArr1 {
		for j,_ := range fArr1[0] {
			arr1[i][j] = int(fArr1[i][j])
		}
	}
	return arr1
} // end func ReduceIntArrayPeak

// ReduceFloat64ArrayAvg reduces arr0 using 'vSliceSize' x 'hSliceSizse' steps
// over arr0; elements of arr1 contain the avg value for a corresponding
// arr0 block; requires power of 2 slice params for valid size reduction
func ReduceFloat64ArrayAvg(arr0 [][]float64, vSliceSize, hSliceSize int) (arr1 [][]float64) {
	lenArrV  := len(arr0)
	lenArrH  := len(arr0[0])
	vSize := lenArrV/vSliceSize // rows in reduced arr1
	hSize := lenArrH/hSliceSize // cols in reduced arr1
	arr1 = make([][]float64, vSize)
	for i,_ := range arr1 {
		arr1[i] = make([]float64, hSize)
	}
	// fmt.Println("--d-- lenArrV:", lenArrV)
	// fmt.Println("--d-- lenArrH:", lenArrH)	
	// fmt.Println("--d-- vSize:", vSize)	
	// fmt.Println("--d-- hSize:", hSize)	
	for i:=0; i<lenArrV; i+=vSliceSize {
		for j:=0; j<lenArrH; j+=hSliceSize {
			// fmt.Printf("--d-- st(%d,%d)\n\r", i,j)
			// --d-- fmt.Println("arr0:", arr0)
			arrSub :=SubsliceFloat64(arr0, i, j, vSliceSize, hSliceSize)
			// fmt.Printf("--d-- sub: %v\n\r", arrSub)
			avg := SliceAvgFloat64(arrSub)
			// --debug-- easier to debug 'sum' array  
			// --debug-- avg := SliceSumFloat64(arrSub) // --debug-- 
			// fmt.Printf("--d-- avg: %v\n\r", avg)
			arr1[i/vSliceSize][j/hSliceSize] = avg
		}
	} // end for i
	return arr1
} // end func ReduceFloat64ArrayAvg

// ReduceFloat64ArrayPeak reduces arr0 using 'vSliceSize' x 'hSliceSizse' steps
// over arr0; elements of arr1 contain the avg value for a corresponding
// arr0 block; requires power of 2 slice params for valid size reduction
func ReduceFloat64ArrayPeak(arr0 [][]float64, vSliceSize, hSliceSize int) (arr1 [][]float64) {
	lenArrV  := len(arr0)
	lenArrH  := len(arr0[0])
	vSize := lenArrV/vSliceSize // rows in reduced arr1
	hSize := lenArrH/hSliceSize // cols in reduced arr1
	arr1 = make([][]float64, vSize)
	for i,_ := range arr1 {
		arr1[i] = make([]float64, hSize)
	}
	// fmt.Println("--d-- lenArrV:", lenArrV)
	// fmt.Println("--d-- lenArrH:", lenArrH)	
	// fmt.Println("--d-- vSize:", vSize)	
	// fmt.Println("--d-- hSize:", hSize)	
	for i:=0; i<lenArrV; i+=vSliceSize {
		for j:=0; j<lenArrH; j+=hSliceSize {
			// fmt.Printf("--d-- st(%d,%d)\n\r", i,j)
			// --d-- fmt.Println("arr0:", arr0)
			arrSub :=SubsliceFloat64(arr0, i, j, vSliceSize, hSliceSize)
			// fmt.Printf("--d-- sub: %v\n\r", arrSub)
			peak := SlicePeakFloat64(arrSub)
			// --debug-- easier to debug 'sum' array  
			// --debug-- peak := SliceSumFloat64(arrSub) // --debug-- 
			// fmt.Printf("--d-- peak: %v\n\r", peak)
			arr1[i/vSliceSize][j/hSliceSize] = peak
		}
	} // end for i
	return arr1
} // end func ReduceFloat64ArrayPeak

// SubsliceFloat64 receives a float array and returns a subset extracting
// nrows, ncols from starting location row, col.
func SubsliceFloat64(
	arr [][]float64, row, col, nrows, ncols int) (arrx [][]float64) {
	arrx = make([][]float64, nrows)
	for i,_ := range arrx {
		arrx[i] = make([]float64, ncols)
	}
	lenArr := len(arr)
	lenArr0 := len(arr[0])
	for i,_ := range arrx {
		for j,_ := range arrx[0] {
			if i+row < lenArr && j+col < lenArr0 {
				arrx[i][j] = arr[i+row][j+col]
			}
		}
	}
	return arrx
}

// SliceAvgFloat64 returns the scalar average of array 'arr0' with
// dimensions 'row' x 'col'
func SliceAvgFloat64(arr0 [][]float64) float64 {
	rowsArr0  := len(arr0)
	colsArr0  := len(arr0[0])
	sum := 0.0
	for i,_ := range arr0 {
		for j:=0; j<colsArr0; j++ {
			sum += arr0[i][j]
		}
	}
	return  sum / float64(rowsArr0*colsArr0) 
} // end func SliceAvgFloat64

// SliceSumFloat64 returns the scalar sum of array 'arr0' with
// dimensions 'row' x 'col'
func SliceSumFloat64(arr0 [][]float64) float64 {
	colsArr0  := len(arr0[0])
	sum := 0.0
	for i,_ := range arr0 {
		for j:=0; j<colsArr0; j++ {
			sum += arr0[i][j]
		}
	}
	return  sum
} // end func SliceSumFloat64

// SliceSumInt returns the scalar sum of array 'arr0' with
// dimensions 'row' x 'col'
func SliceSumInt(arr0 [][]int) int {
	colsArr0  := len(arr0[0])
	sum := 0
	for i,_ := range arr0 {
		for j:=0; j<colsArr0; j++ {
			sum += arr0[i][j]
		}
	}
	return  sum
} // end func SliceSumInt

// SlicePeakFloat64 returns the scalar peak value of array 'arr0' with
// dimensions 'row' x 'col'
func SlicePeakFloat64(arr0 [][]float64) float64 {
	colsArr0  := len(arr0[0])
	peak := arr0[0][0]
	for i,_ := range arr0 {
		for j:=0; j<colsArr0; j++ {
			if arr0[i][j] > peak {
				peak = arr0[i][j]
			}
		}
	}
	return  peak
} // end func SlicePeakFloat64



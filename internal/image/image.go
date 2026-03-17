package image

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"

	"github.com/spectriclabs/sigplot-data-service/internal/util"
)

func ApplyCXmode(datain []float64, cxmode string, complexData bool) []float64 {
	loThresh := 1.0e-20
	if complexData {
		outData := make([]float64, len(datain)/2)
		for i := 0; i < len(datain)-1; i += 2 {
			switch cxmode {
			case "Ma":
				outData[i/2] = math.Sqrt(datain[i]*datain[i] + datain[i+1]*datain[i+1])
			case "Ph":
				outData[i/2] = math.Atan2(datain[i+1], datain[i])
			case "Re":
				outData[i/2] = datain[i]
			case "Im":
				outData[i/2] = datain[i+1]
			case "IR":
				outData[i/2] = math.Sqrt(datain[i]*datain[i] + datain[i+1]*datain[i+1])
			case "Lo":
				mag2 := datain[i]*datain[i] + datain[i+1]*datain[i+1]
				mag2 = math.Max(mag2, loThresh)
				outData[i/2] = 10 * math.Log10(mag2)
			case "L2":
				mag2 := datain[i]*datain[i] + datain[i+1]*datain[i+1]
				mag2 = math.Max(mag2, loThresh)
				outData[i/2] = 20 * math.Log10(mag2)
			default: //Defaults to "real"
				log.Println("Unkown cxmode", cxmode, "defaulting to Re")
				outData[i/2] = datain[i]
			}
		}
		return outData
	} else {
		switch cxmode {
		case "Ma":
			outData := make([]float64, len(datain))
			for i := range datain {
				outData[i] = math.Abs(datain[i])
			}
			return outData
		case "Ph":
			outData := make([]float64, len(datain))
			for i := 0; i < len(datain); i++ {
				outData[i] = math.Atan2(0, datain[i])
			}
			return outData
		case "Re":
			return datain
		case "Im":
			outData := make([]float64, len(datain))
			return outData
		case "IR":
			return datain
		case "Lo":
			for i := range datain {
				datain[i] = 10 * math.Log10(math.Max(datain[i], loThresh))
			}
			return datain
		case "L2":
			for i := range datain {
				datain[i] = 20 * math.Log10(math.Max(datain[i], loThresh))
			}
			return datain

		}
		return datain //Defaults to "Real" or passthrough

	}
}

func DownSampleLineInY(datain []float64, outxsize int, transform string) []float64 {
	numLines := len(datain) / outxsize
	//log.Println("len(datain),outxsize" ,len(datain),outxsize)
	processSlice := make([]float64, numLines)
	outData := make([]float64, outxsize)
	for x := 0; x < outxsize; x++ {
		for y := 0; y < numLines; y++ {
			//log.Println("y thin" ,y,outxsize,x)
			processSlice[y] = datain[y*outxsize+x]
		}
		outData[x] = Transform(processSlice[:], transform)
	}
	return outData
}

func DownSampleLineInX(datain []float64, outxsize int, transform string, outData []float64, outLineNum int) {
	//var inputysize int =len(datain)/framesize
	var xelementsperoutput float64
	xelementsperoutput = float64(len(datain)) / float64(outxsize)
	//var thinxdata = make([]float64,outxsize)
	if xelementsperoutput > 1 { // Expansion
		for x := 0; x < outxsize; x++ {
			var startelement int
			var endelement int
			if x != (outxsize - 1) { // Not last element
				startelement = int(math.Round(float64(x) * xelementsperoutput))
				endelement = int(math.Round(float64(x+1) * xelementsperoutput))
			} else { // Last element, work backwards
				endelement = len(datain)
				startelement = endelement - int(math.Ceil(xelementsperoutput))
			}

			outData[outLineNum*outxsize+x] = Transform(datain[startelement:endelement], transform)

		}
	} else { // Expand Data by repeating input values into output

		for x := 0; x < outxsize; x++ {
			index := int(math.Floor(float64(x) * xelementsperoutput))
			outData[outLineNum*outxsize+x] = datain[index]
		}
	}
}

func Transform(dataIn []float64, transform string) float64 {
	switch transform {
	case "mean":
		var sum float64
		for _, v := range dataIn {
			sum += v
		}
		num := sum / float64(len(dataIn))
		if math.IsNaN(num) {
			log.Println("DoTransform produced NaN")
			num = 0
		}
		return num
	case "max":
		num := dataIn[0]
		for _, v := range dataIn[1:] {
			if v > num {
				num = v
			}
		}
		if math.IsNaN(num) {
			log.Println("DoTransform produced NaN")
			num = 0
		}
		return num
	case "min":
		num := dataIn[0]
		for _, v := range dataIn[1:] {
			if v < num {
				num = v
			}
		}
		if math.IsNaN(num) {
			log.Println("DoTransform produced NaN")
			num = 0
		}
		return num
	case "maxabs":
		maxVal := math.Abs(dataIn[0])
		for _, v := range dataIn[1:] {
			if av := math.Abs(v); av > maxVal {
				maxVal = av
			}
		}
		if math.IsNaN(maxVal) {
			log.Println("DoTransform produced NaN")
			maxVal = 0
		}
		return maxVal
	case "first":
		num := dataIn[0]
		if math.IsNaN(num) {
			log.Println("DoTransform produced NaN")
			num = 0
		}
		return num
	default: // Default to first if bad value.
		log.Println("Unknown transform", transform, "using first")
		num := dataIn[0]
		if math.IsNaN(num) {
			log.Println("DoTransform produced NaN")
			num = 0
		}
		return num

	}
}

func CreateOutput(dataIn []float64, fileFormat string, zmin, zmax float64, colorMap string) []byte {
	// for i := 0; i < len(dataIn); i++ {
	// 	if math.IsNaN(dataIn[i]) {
	// 		log.Println("CreateOutput NaN", i)
	// 	}
	// }

	dataOut := new(bytes.Buffer)
	numColors := 1000
	if fileFormat == "RGBA" {
		colorPalette := GetCachedPalette(colorMap, numColors)
		if zmax != zmin {
			colorsPerSpan := (zmax - zmin) / float64(numColors)
			output := make([]byte, len(dataIn)*4)
			for i, v := range dataIn {
				colorIndex := math.Round((v-zmin)/colorsPerSpan) - 1
				colorIndex = math.Min(math.Max(colorIndex, 0), float64(numColors-1))
				ci := int(colorIndex)
				offset := i * 4
				output[offset] = byte(colorPalette[ci].Red)
				output[offset+1] = byte(colorPalette[ci].Green)
				output[offset+2] = byte(colorPalette[ci].Blue)
				output[offset+3] = 255
			}
			return output
		} else {
			output := make([]byte, len(dataIn)*4)
			r := byte(colorPalette[0].Red)
			g := byte(colorPalette[0].Green)
			b := byte(colorPalette[0].Blue)
			for i := 0; i < len(dataIn); i++ {
				offset := i * 4
				output[offset] = r
				output[offset+1] = g
				output[offset+2] = b
				output[offset+3] = 255
			}
			return output
		}
	} else {
		log.Println("Creating Output of Type ", fileFormat)
		switch string(fileFormat[1]) {
		case "B":
			var numSlice = make([]int8, len(dataIn))
			for i := 0; i < len(numSlice); i++ {
				numSlice[i] = int8(math.Round(dataIn[i]))
			}

			err := binary.Write(dataOut, binary.LittleEndian, &numSlice)

			util.CheckError(err)

		case "I":
			var numSlice = make([]int16, len(dataIn))
			for i := 0; i < len(numSlice); i++ {
				numSlice[i] = int16(math.Round(dataIn[i]))
			}

			err := binary.Write(dataOut, binary.LittleEndian, &numSlice)

			util.CheckError(err)

		case "L":
			var numSlice = make([]int32, len(dataIn))
			for i := 0; i < len(numSlice); i++ {
				numSlice[i] = int32(math.Round(dataIn[i]))
			}

			err := binary.Write(dataOut, binary.LittleEndian, &numSlice)

			util.CheckError(err)

		case "F":
			var numSlice = make([]float32, len(dataIn))
			for i := 0; i < len(numSlice); i++ {
				numSlice[i] = float32(dataIn[i])
			}

			err := binary.Write(dataOut, binary.LittleEndian, &numSlice)

			util.CheckError(err)

		case "D":
			var numSlice = make([]float64, len(dataIn))
			for i := 0; i < len(numSlice); i++ {
				numSlice[i] = dataIn[i]
			}

			err := binary.Write(dataOut, binary.LittleEndian, &numSlice)

			util.CheckError(err)

		case "P":
			extraBits := len(dataIn) % 8
			for extraBit := 0; extraBit < extraBits; extraBits++ { //Pad zeros to make the number of elements divisable by 8 so it can be packed into a byte
				dataIn = append(dataIn, 0)
			}
			numBytes := len(dataIn) / 8
			var numSlice = make([]uint8, numBytes)
			for i := 0; i < len(numSlice); i++ {
				for j := 0; j < 8; j++ {
					var bit uint8
					if dataIn[i*8+j] > 0 { //SP Data can only be 0 or 1, so if values is greater than 0, make it a 1.
						bit = 1
					} else {
						bit = 0
					}
					numSlice[i] = (numSlice[i] << 1) | bit
				}

			}
			err := binary.Write(dataOut, binary.LittleEndian, &numSlice)
			util.CheckError(err)

		default:
			log.Println("Unsupported output type")
		}
		//log.Println("out_data" , len(dataOut.Bytes()))

		return dataOut.Bytes()
	}

}

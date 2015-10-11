package main

import
(	"fmt"
	"flag"
	"image"
	"image/png"
	"os"
	"encoding/binary"
	_ "strconv"
	"math"
)

type Pixel struct {
	x int
	y int
	rgba string
}


func main() {
	fmt.Println("LSB Goooo")	
	var container, message, output string
	var extract bool
	flag.StringVar(&container, "container", "", "The container input")
	flag.StringVar(&message, "message", "", "Message file")
	flag.StringVar(&output, "output", "", "The container output")
	flag.BoolVar(&extract, "extract", false, "If set : extraction")
	flag.Parse()

	var error bool = false

	/*********************/
	//Check params
	/*********************/
	//global
	if len(container)==0 {
		fmt.Println("Please set 'container' params for insertion/extraction.")
		error = true
	}

	//extraction
	if extract {
		if len(output)==0 {
			fmt.Println("Please set 'output' params for extraction.")
			error = true
		}	
	}else
	//insertion
	{
		if len(message)==0 {
			fmt.Println("Please set 'message' params for insertion.")
			error = true	
		}
		if len(output)==0 {
			fmt.Println("Please set 'output' params for insertion.")
			error = true
		}	
	}

	if error {
		fmt.Println("There is errors, please correct them.")
		return
	}



	//action
	if extract {
		extractMessage(container, output)
	}else{
		insertMessage(container,message,output)
	}

}


func insertMessage(containerPath string, messagePath string, outputPath string){
	container := openImgPng(containerPath)
	bounds := container.Bounds()

	messageBytes := openMessage(messagePath)

	pixel := Pixel{-1,-1,""}

	lenData := len(messageBytes)
	tmpLen := lenData
	for i:=0;i<32;i++ {
		bytLSB := getPixelBytes(container,pixel)
		
		//get bit of message
		bit := tmpLen%2
		tmpLen=tmpLen/2

		//set bit on container
		bytLSB = setLSB(bytLSB,bit)
		setPixelBytes(container, pixel, bytLSB)

		pixel.next(bounds)
	}

	
	for i:=0;i<lenData;i++ {
		//fmt.Println(messageBytes[i])
		for j:=0;j<8;j++ {
			//get byte of container
			bytLSB := getPixelBytes(container,pixel)
			
			//get bit of message
			bit := getMessageBit(messageBytes,i,j)
			// fmt.Print(bit)

			//set bit on container
			bytLSB = setLSB(bytLSB,bit)
			setPixelBytes(container, pixel, bytLSB)


			pixel.next(bounds)
		}
		// fmt.Println("")
				
	}

	outimg, err := os.Create(outputPath)
	check(err)
	defer outimg.Close()
	err = png.Encode(outimg, container.SubImage(container.Bounds()))
	check(err)
}

func extractMessage(containerPath string, outputPath string){
	container := openImgPng(containerPath)
	bounds := container.Bounds()
	pixel := Pixel{-1,-1,""}
	length:=0
	for i:=0;i<32;i++ {
		bytLSB := getPixelBytes(container,pixel)
		length += getLSB(bytLSB)*int(math.Pow(float64(2),float64(i)))
		
		pixel.next(bounds)
	}

	tab := make([]byte, int(length))
	var by byte

	for i:=0;i<length;i++ {

		for j:=0;j<8;j++ {
			
			
			bytLSB := getPixelBytes(container,pixel)
			bit := getLSB(bytLSB)
			// fmt.Print(bit)
			
			by = setMessageBit(by,j,bit)


			pixel.next(bounds)
		}
		
		tab[i]=by
		// fmt.Println("")
		// fmt.Println(by)
		by = 0	
		
	}
	// fmt.Println("")
	// fmt.Println(by)
	// fmt.Println(tab)
	writeMessage(outputPath,tab)
}


func (p *Pixel) next(b image.Rectangle){
	if p.x==-1 {
		p.x=0
		p.y=0
		p.rgba="r" 
		return
	}
	//TODO seek alpha

	/*switch p.rgba {
		case "r":
			p.rgba = "g"
		case "g":
			p.rgba = "b"
		case "b":
			p.rgba = "r"
	}*/

	if p.rgba=="r" {
		if p.x < b.Max.X {
			p.x++
		}else{
			p.x=0
			p.y++
			if p.y>b.Max.Y {
				fmt.Println("TODO check size before.")
			}
		}
	}
}

func getMessageBit(b []byte, numByte int, numBit int)(int){
	return ( int(b[numByte]) & (1<<uint(numBit-1))) >> uint(numBit-1)
}

func setMessageBit(b byte, numBit int, value int)(byte){
	v1 := int(b)
	v2 := v1 &^ ((1<<uint(numBit-1)))
	v3 := v2 | (value<<uint(numBit-1))
	return byte(v3)

}

func setLSB(b byte, value int)(byte){
	return  byte ( ( (int(b) >> uint(1) ) << uint(1)) + value)
}

func getLSB(b byte)(int){
	return int(b)%2
}

func getPixelBytes(img image.NRGBA, pixel Pixel)(byte){
	rgba := img.NRGBAAt(pixel.x, pixel.y)
	switch pixel.rgba {
		case "r" :
			return rgba.R
		case "g" :
			return rgba.G
		case "b" :
			return rgba.B
		case "a" :
			return rgba.A
	}
	return 0
}

func setPixelBytes(img image.NRGBA, pixel Pixel, byt byte){
	rgba := img.NRGBAAt(pixel.x, pixel.y)
	switch pixel.rgba {
		case "r" :
			rgba.R = byt
		case "g" :
			rgba.G = byt
		case "b" :
			rgba.B = byt
		case "a" :
			rgba.A = byt
	}
	img.Set(pixel.x, pixel.y, rgba)
}

func openMessage(filename string)([]byte){
	f, err := os.Open(filename)
	if err != nil {
	   panic(err)
	}
	defer f.Close()	
	stats, err := f.Stat()
	check(err)
	
	tab := make([]byte, stats.Size())
	_, err = f.Read(tab)
	check(err)

	return tab
}

func writeMessage(filename string, tab []byte){
	f, err := os.Create(filename)
	if err != nil {
	   panic(err)
	}
	defer f.Close()	
	
	
	
	n2, err := f.Write(tab)
    check(err)
    fmt.Printf("wrote %d bytes\n", n2)
}

func openImgPng(filename string) (image.NRGBA){

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Path error image")
		os.Exit(-1)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("Error in image")
		os.Exit(-1)
	}

	im := img.(*image.NRGBA)
	return *im
}


func encodeLen(l int) ([]byte){
	/*tab := make([]byte, 8)
	i:=0
	j:=0
	for l!=0 && i<8 {
		tab[j] = (tab[j]<<1) + l%2
		
		l=l/2
		i++
		j=i/8
	}
	return tab*/
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Println(l)
	bs := make([]byte, 4)
    binary.LittleEndian.PutUint32(bs, uint32(l))
    fmt.Println(bs)
    for i:=0;i<32;i++ {
	    bit := getMessageBit(bs,int(i/8),i%8)
	    fmt.Print(bit)
    }
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~")
    return bs
}

func decodeLen(tab []byte)(int){
	fmt.Println("************************")
	fmt.Println(tab)
	nb :=int(binary.LittleEndian.Uint32(tab))
	fmt.Println(nb)	
	fmt.Println("************************")
	return nb
}



func check(e error) {
    if e != nil {
        panic(e)
    }
}
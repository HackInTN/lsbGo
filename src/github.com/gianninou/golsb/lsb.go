package main

import
(	
	"fmt"
	"flag"
	"image"
	"image/png"
	"os"
	_ "strconv"
	"math"


)


var container, message, output string
var extract bool

var pixelPath, rgbPath string
var key string


func main() {	


	flag.StringVar(&container, "container", "", "The container input")
	flag.StringVar(&message, "message", "", "Message file")
	flag.StringVar(&output, "output", "", "The container output")
	flag.BoolVar(&extract, "extract", false, "If set : extraction")

	flag.StringVar(&pixelPath, "pixelPath", "horizontal", "Pixel Path methode (horizontal/vertical/diagonal/field)")

	flag.StringVar(&rgbPath, "rgbPath", "basic", "RGB Path methode (none/basic)")

	flag.StringVar(&key, "key", "", "Crypt/decrypt key")


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
	//initialisation
	container := openImgPng(containerPath)
	bounds := container.Bounds()

	maxSize := bounds.Max.X*bounds.Max.Y
	if rgbPath=="basic" {
		maxSize *= 3
	}
	fmt.Println("Maximum message size ",maxSize)

	messageBytes := openMessage(messagePath)
	cypherMessageBytes := encrypt(messageBytes,[]byte(key))

	

	pixel := IteratorInit()
	pixel.SetLayoutPath(pixelPath)
	pixel.SetRgbPath(rgbPath)
	pixel.PrintPath()
	nbWrite:=0

	lenData := len(cypherMessageBytes)
	fmt.Println("Message size : ",lenData)
	tmpLen := lenData
	for i:=0;i<32;i++ {
		bytLSB := getPixelBytes(container,pixel)
		
		//get bit of message
		bit := tmpLen%2
		
		tmpLen=tmpLen/2

		//set bit on container
		bytLSB = setLSB(bytLSB,bit)
		setPixelBytes(container, pixel, bytLSB)

		pixel.Next(bounds.Max.X,bounds.Max.Y)
	}

	var by byte
	for i:=0;i<lenData;i++ {
		by = cypherMessageBytes[i]
		for j:=0;j<8;j++ {
			//get byte of container
			bytLSB := getPixelBytes(container,pixel)
			
			//get bit of message
			bit := getMessageBit(by,j)

			//set bit on container
			bytLSB = setLSB(bytLSB,int(bit))
			setPixelBytes(container, pixel, bytLSB)
			nbWrite++

			pixel.Next(bounds.Max.X,bounds.Max.Y)
		}
				
	}

	outimg, err := os.Create(outputPath)
	check(err)
	defer outimg.Close()
	err = png.Encode(outimg, container.SubImage(container.Bounds()))
	check(err)
	fmt.Printf("Ecriture de %d octets\n",(nbWrite/8))
}

func extractMessage(containerPath string, outputPath string){
	container := openImgPng(containerPath)
	bounds := container.Bounds()
	pixel := IteratorInit()
	pixel.SetLayoutPath(pixelPath)
	pixel.SetRgbPath(rgbPath)

	length:=0
	for i:=0;i<32;i++ {
		bytLSB := getPixelBytes(container,pixel)
		length += getLSB(bytLSB)*int(math.Pow(float64(2),float64(i)))
		
		pixel.Next(bounds.Max.X,bounds.Max.Y)
	}

	tab := make([]byte, uint(length))
	var by byte

	for i:=0;i<length;i++ {

		for j:=0;j<8;j++ {
			
			
			bytLSB := getPixelBytes(container,pixel)
			bit := getLSB(bytLSB)
			
			by = setMessageBit(by,j,bit)


			pixel.Next(bounds.Max.X,bounds.Max.Y)
		}
		
		tab[i]=by
		by = 0x0000	
		
	}
	
	decryptTab := decrypt(tab,[]byte(key))
	writeMessage(outputPath,decryptTab)
}




func getMessageBit(b byte,  numBit int)(uint){
	return ( uint(b) & (1<<uint(numBit)) ) >> uint(numBit)
}

func setMessageBit(b byte, numBit int, value int)(byte){
	v1 := int(b)
	v2 := v1 &^ ((1<<uint(numBit)))
	v3 := v2 | (value<<uint(numBit))
	return byte(v3)

}

func setLSB(b byte, value int)(byte){
	return  byte ( ( (int(b) >> uint(1) ) << uint(1)) + value)
}

func getLSB(b byte)(int){
	return int(b)%2
}

func getPixelBytes(img image.NRGBA, pixel LsbPixel)(byte){
	rgba := img.NRGBAAt(pixel.GetX(), pixel.GetY())
	switch pixel.GetLayer() {
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

func setPixelBytes(img image.NRGBA, pixel LsbPixel, byt byte){
	rgba := img.NRGBAAt(pixel.GetX(), pixel.GetY())
	switch pixel.GetLayer() {
		case "r" :
			rgba.R = byt
		case "g" :
			rgba.G = byt
		case "b" :
			rgba.B = byt
		case "a" :
			rgba.A = byt
	}
	img.Set(pixel.GetX(), pixel.GetY(), rgba)
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
    fmt.Printf("Extraction de %d octets\n", n2)
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




func check(e error) {
    if e != nil {
        panic(e)
    }
}

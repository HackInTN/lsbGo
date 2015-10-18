package main

import
(	
	"fmt"
	"os"
)

const (
	HORIZONTAL=iota
	VERTICAL=iota
	DIAGONAL=iota
	FIELD=iota
)


var TAB_PIXEL_CONST = []string{"horizontal","vertical","diagonal","field"}
var	TAB_RGB_CONST = []string{"none","basic"}


const (
	NONE=iota
	BASIC=iota
)

type LsbPixel struct {
	x int
	y int
	rgba string
	methode int
	layer int
}

func IteratorInit() LsbPixel{
	return LsbPixel{0,0,"r",HORIZONTAL,BASIC}
}


func (p *LsbPixel) SetLayoutPath(path string){
	switch path {
		case "horizontal":
			p.methode=HORIZONTAL
		case "vertical":
			p.methode=VERTICAL
		case "diagonal":
			p.methode=DIAGONAL
		case "field":
			p.methode=FIELD
	}
}

func (p *LsbPixel) SetRgbPath(path string){
	switch path {
		case "none":
			p.layer=NONE
		case "basic":
			p.layer=BASIC
	}
}



func (p *LsbPixel) GetX() int{
	return p.x
}

func (p *LsbPixel) GetY() int{
	return p.y
}

func (p *LsbPixel) GetLayer() string{
	return p.rgba
}

func (p *LsbPixel) Next(dx, dy int){
	
	next := p.NextLayer()
	
	if next {
		switch p.methode {
			case HORIZONTAL:
				p.NextHorizontal(dx,dy)
			case VERTICAL:
				p.NextVertical(dx,dy)
			case DIAGONAL:
				p.NextDiagonal(dx,dy)
			case FIELD:
				p.NextCorps(dx,dy)
		}
	}
	//fmt.Println(p)
}


func (p *LsbPixel) NextLayer()bool{
	var res bool
	switch p.layer {
		case NONE:
			res = true
		case BASIC:
			res = p.NextLayerBasic()
	}
	return res
}

func (p *LsbPixel) NextLayerBasic()bool{
	switch p.rgba {
		case "r":
			p.rgba = "g"
			return false
		case "g":
			p.rgba = "b"
			return false
		case "b":
			p.rgba = "r"
			return true
	}
	return false
}


func (p *LsbPixel) NextHorizontal(dx, dy int){
	if p.rgba=="r" {
		p.x++
		if p.x == dx {
			p.x=0
			p.y++
			if p.y>=dy {
				fmt.Println("TODO check size before.")
			}
		}
	}
}


func (p *LsbPixel) NextVertical(dx, dy int){
	if p.rgba=="r" {
		p.y++
		if p.y == dy {
			p.y=0
			p.x++
			if p.x>=dx {
				fmt.Println("TODO check size before.")
			}
		}
	}
}

func (p *LsbPixel) NextDiagonal(dx, dy int){
	curr := p.x+p.y

	p.x++
	p.y--

	for p.x>=dx || p.y>=dy || p.x<0 || p.y<0 {

		if p.y<0 || p.x>=dx {
			curr++
			if curr>dx+dy {
				fmt.Println("Erreur borne : Exit")
				os.Exit(-1)
			}
			p.x=0
			p.y=curr
		}else{
			p.x++
			p.y--			
		}
	}
}


func (p *LsbPixel) NextCorps(dx, dy int){
	//fmt.Println(dx,dy)
	if p.y==0 && p.x==0 {
		p.x=1
		return
	}

	N := 65537
	P := 3

	nb := p.y*dx + p.x
	nb = (nb*P)%N

	for nb/dx>=dy || nb%dx>=dx {
		nb = (nb*P)%N
	}
	p.x=nb%dx
	p.y=nb/dx

	if p.y==0 && p.x==1 {
		fmt.Println("Erreur borne : Exit")
		os.Exit(-1)
	}
	//fmt.Println(p)

}

func (p *LsbPixel) PrintPath(){
	fmt.Println("Méthode : ",TAB_PIXEL_CONST[p.methode],", RGB : ",TAB_RGB_CONST[p.layer])
}
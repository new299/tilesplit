package main

import (
  "fmt"
  "image"
  "image/jpeg"
  "image/png"
  "image/draw"
  "os"
  "strconv"
  "github.com/nfnt/resize"
  "math"
)

func init() {
  // damn important or else At(), Bounds() functions will
  // caused memory pointer error!!
  image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
  image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

func make_image_tiles(basepath string,scale int,tile_size int,rescale_size int,rgbimage *image.RGBA) {
  bounds   := rgbimage.Bounds()
  fmt.Println(bounds)

  path := basepath + "/" + strconv.Itoa(scale)
  os.Mkdir(path,777);

  for cx := bounds.Min.X;cx < bounds.Max.X;cx += tile_size {

  os.Mkdir(path + "/" + strconv.Itoa(cx/tile_size),777);

   for cy := bounds.Min.Y;cy < bounds.Max.Y;cy += tile_size {
  
     // Grab tile from image
     fmt.Printf("Get tile %v %v %v %v\n",cx,cy,cx+tile_size,cy+tile_size)
     subimage := rgbimage.SubImage(image.Rectangle{image.Point{cx,cy},image.Point{cx+tile_size,cy+tile_size}})

     subbounds := subimage.Bounds()
     fmt.Println("subb: ", subbounds)
     x_delta := (subbounds.Max.X-subbounds.Min.X)
     y_delta := (subbounds.Max.Y-subbounds.Min.Y)
     fmt.Println("delta: ", x_delta," ",y_delta)
     if (x_delta < tile_size) ||
        (y_delta < tile_size) {

       newsubimage := image.NewRGBA(image.Rectangle{image.Point{0,0},image.Point{tile_size,tile_size}})
       draw.Draw(newsubimage,image.Rectangle{image.Point{0,0},image.Point{tile_size,tile_size}},subimage,subimage.Bounds().Min,draw.Src)
       subimage = newsubimage

     }
     fmt.Println("subdixed: ", subimage.Bounds())

     if tile_size != rescale_size {
      subimage = resize.Resize(uint(rescale_size),uint(rescale_size), subimage, resize.Lanczos3)
     }

     // Write the file to disk
     subfile, _ := os.Create(path + "/" + strconv.Itoa(cx/tile_size) + "/" + strconv.Itoa(cy/tile_size) + ".png")
     png.Encode(subfile, subimage) 
   }
  }
}

func main() {
  imgfile, err := os.Open(os.Args[1])

  tile_size := 256
  rescale_size := tile_size

  if err != nil {
    fmt.Println("file not found!")
    os.Exit(1)
  }

  defer imgfile.Close()

  img, _, _ := image.Decode(imgfile)

  rgbimage := img.(*image.RGBA)
  bounds   := rgbimage.Bounds()

  width  := bounds.Max.X - bounds.Min.X
  height := bounds.Max.Y - bounds.Max.Y

  // there's no int math.Max in golang
  mdim := math.Max(float64(width),float64(height))

  scale_max_f := math.Log2(mdim/float64(tile_size))+1
  scale_max   := int(scale_max_f)

  os.Mkdir("./tiles",777);

  c_tile_size := tile_size
  for scale := scale_max; scale >= 1; scale-- {
    fmt.Printf("tile size: %v\n",c_tile_size)
    make_image_tiles("./tiles",scale,int(c_tile_size),rescale_size,rgbimage);
    c_tile_size = c_tile_size*2
  }

}

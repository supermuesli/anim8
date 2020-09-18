# anim8
---
*anim8* is a simple animation/storyboarding tool written in go. You will need *GLFW* in order to run this.

## Showcase
![alt text](https://github.com/supermuesli/anim8/blob/master/showcase.gif "y u trippin")

## Build
Simply do
```  
$ make
```
or if *make* is not an option for you, just do:
``` 
$ go build -o $(GOPATH)/bin/anim8 cmd/anim8/main.go cmd/anim8/bindata.go
```
and make sure *$(GOPATH)/bin* is in your *PATH*. To run *anim8*, you should now be able to do
``` 
$ anim8
```

## Usage
- start sketching your first frame
- now press **SPACE**, it will store the frame in a buffer (which will become the animation)
- you will notice that the previous frame is still showing with 30% opacity, as a guide for the next frame
- continue collecting frames until you think you have enough
- if you made a mistake and need to erase from the canvas using your cursor, press **E**
- if you want to delete the entire frame, press **D**
- if you need to adjust the brush size, use your **mouse wheels** to do so
- if you want to see how your frames look animated, press **P**
- if you want to dump your animation as a set of PNGs, press **ENTER**, type the scene name, and type **ENTER** again
- if you want to start sketching the next scene, press **R** and start collecting frames again
- that's pretty much the intended workflow

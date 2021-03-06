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
- now press **SPACE**, it will store the frame in a scene (a buffer which will become the animation)
- you will notice that the previous frame is still showing with 30% opacity, as a guide for the next frame (the indication won't be stored in the scene)
- you can also press **C** *(copy)* to copy (and overwrite) the previous frame to the current one 
- using **SHIFT** + arrow keys you can shift the current frame in a certain direction
- you can press **E** *(erase)* to erase from the canvas using your cursor, then press **E** again to return to the brush tool 
- if you need to update an older frame, you can use the **LEFT** and **RIGHT** arrow keys to nagivate through the scene
- if you want to delete the entire frame, press **D** *(delete)*
- if you need to adjust the brush size, use your **mouse wheels** to do so
- continue collecting frames until you think you have enough
- if you want to see how your frames look animated, press **P** *(play)*
  - you can cancel the animation by pressing and holding **P** again
- if you want to view the animation in a loop, press **L** *(loop)*, and to escape the loop press and hold **L** again
  - while in loop-mode, you can press and hold the **UP** and **DOWN** arrow keys to increase or decrease the playback FPS
- if you want to dump your animation as a set of PNGs, press **ENTER**, type the scene name, and press **ENTER** again
- if you want to reset the scene and start collecting frames for another, press **R** *(reset)*
- that's pretty much the intended workflow
- press **ESC** *(escape)* to exit the program
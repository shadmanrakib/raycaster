# Raycaster

https://github.com/shadmanrakib/raycaster/assets/64807913/fbfbb9ef-68d9-4f73-b304-3674a4cb7a89


I really wanted to explore raycasting and some of the math behind it. 
So, I created this raycasting sandbox to experiment with it. 
On the left half of the screen, I have a top-down 2D view of the
map, with rays being casted from the player. On the right half, I
implemented a 3D view of the scene using the euclidean distances I
got from raycasting. I've left the fisheye effect (distortion) 
because I thought it looked cool. It's possible to correct the 
distortion by using the perpendicular distance or make sure to cast
rays into the pixels of the camera plane instead of evenly spacing them
out.

## Movement and Keybindings

Movement is relative to the 2D top down view. To move around in the 
demo, use the WASD keys or the arrow keys. To rotate clockwise, use 
the F key or the Space key. To rotate counter-clockwise use the G key.

## Installation

Please install Golang 1.21 or higher. After doing so, either do
`go run main.go` or `go build main.go && ./main`. The latter 
compiles the code into an executable and runs it. Note the above
instructions should work on most MacOS and Linux machines.

## Learnings

Using this much trignometry was quite fun yet stressful. 
It might be much easier to define the camera plane and rely less on
trignometry and use the fact that a vector has two components, and
for calculating the distance, we really only need one of these
components if we know how much the length of the ray increases
whenever one of the components increases by 1. This might help
reduce floating point errors, which were quite dissatisfying. 


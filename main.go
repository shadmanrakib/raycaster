package main

import (
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/shadmanrakib/raycast/utils"
)

var MIN_ROTATION float64 = 0
var MAX_ROTATION float64 = 2 * math.Pi

type Camera struct {
	// xy_rotation is in rads
	x, y, xy_rotation, fov float64
}

type Game struct {
	width, height, scale_factor                                           int
	grid                                                                  [][]bool
	camera                                                                Camera
	player_size                                                           float64
	key_down_is_down, key_up_is_down, key_left_is_down, key_right_is_down bool
	key_rotate_clockwise_is_down, key_rotate_counter_clockwise_is_down    bool
}

func (g *Game) Update() error {
	// movement
	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.key_up_is_down = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.key_down_is_down = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.key_left_is_down = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.key_right_is_down = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyW) || inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		g.key_up_is_down = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyS) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		g.key_down_is_down = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyA) || inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		g.key_left_is_down = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyD) || inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		g.key_right_is_down = false
	}

	if g.key_up_is_down && utils.CanContinueInDirection(g.grid, g.camera.x-g.player_size/2, g.camera.y, -0.05, 0) {
		g.camera.x -= 0.05
	}
	if g.key_down_is_down && utils.CanContinueInDirection(g.grid, g.camera.x+g.player_size/2, g.camera.y, 0.05, 0) {
		g.camera.x += 0.05
	}
	if g.key_left_is_down && utils.CanContinueInDirection(g.grid, g.camera.x, g.camera.y-g.player_size/2, 0, -0.05) {
		g.camera.y -= 0.05
	}
	if g.key_right_is_down && utils.CanContinueInDirection(g.grid, g.camera.x, g.camera.y+g.player_size/2, 0, 0.05) {
		g.camera.y += 0.05
	}

	// rotation
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.key_rotate_clockwise_is_down = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.key_rotate_counter_clockwise_is_down = true
	}

	if inpututil.IsKeyJustReleased(ebiten.KeySpace) || inpututil.IsKeyJustReleased(ebiten.KeyF) {
		g.key_rotate_clockwise_is_down = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyG) {
		g.key_rotate_counter_clockwise_is_down = false
	}

	if g.key_rotate_counter_clockwise_is_down {
		g.camera.xy_rotation -= 0.01
	}

	if g.key_rotate_clockwise_is_down {
		g.camera.xy_rotation += 0.01
	}

	return nil
}

func drawMap2D(g *Game, screen *ebiten.Image) {
	for i := 0; i < len(g.grid); i++ {
		for j := 0; j < len(g.grid[0]); j++ {
			var c color.RGBA
			void_color := color.RGBA{30, 30, 30, 255}
			wall_color := color.RGBA{200, 30, 30, 255}
			if g.grid[i][j] {
				c = wall_color
			} else {
				c = void_color
			}

			vector.DrawFilledRect(screen, float32(j*g.scale_factor), float32(i*g.scale_factor), float32(g.scale_factor)-1, float32(g.scale_factor)-1, c, false)
		}
	}
}

func drawPlayer2D(g *Game, screen *ebiten.Image) {
	c := color.RGBA{200, 200, 30, 255}
	vector.DrawFilledCircle(screen, float32(g.camera.y)*float32(g.scale_factor), float32(g.camera.x)*float32(g.scale_factor), float32(g.player_size)*float32(g.scale_factor)/2, c, false)
}

func drawIntersection2D(g *Game, screen *ebiten.Image, x, y float64) {
	c := color.RGBA{30, 200, 30, 255}
	vector.DrawFilledCircle(screen, float32(y)*float32(g.scale_factor), float32(x)*float32(g.scale_factor), 2, c, false)
}

func draw2D(g *Game, screen *ebiten.Image) {
	drawMap2D(g, screen)
	drawPlayer2D(g, screen)

	start_ray_rad := g.camera.xy_rotation - (g.camera.fov / 2)

	num_rays := 40
	ray_rad_increment := g.camera.fov / (float64(num_rays) - 1)
	ray := start_ray_rad

	for i := 0; i < num_rays; i++ {
		dist, _, _ := utils.DDA(g.grid, g.camera.x, g.camera.y, ray, false)
		norm_x, norm_y := utils.CalcNormDirVectorFromRadians(ray)
		f_x, f_y := dist*norm_x*float64(g.scale_factor)+g.camera.x*float64(g.scale_factor), dist*norm_y*float64(g.scale_factor)+g.camera.y*float64(g.scale_factor)
		vector.StrokeLine(screen, float32(g.camera.y)*float32(g.scale_factor), float32(g.camera.x)*float32(g.scale_factor), float32(f_y), float32(f_x), 1, color.White, false)
		ray += ray_rad_increment
	}
}

func draw3D(g *Game, screen *ebiten.Image, canvas_offset float32) {
	start_ray_rad := g.camera.xy_rotation - (g.camera.fov / 2)
	ray_rad_increment := g.camera.fov / (float64(g.width) - 1)
	h := g.height
	ray := start_ray_rad

	for i := 0; i < g.width; i++ {
		dist, _, _ := utils.DDA(g.grid, g.camera.x, g.camera.y, ray, false)

		// Height of line to draw
		lineHeight := (int)(float64(h) / dist)

		// Find start and end of line, cap the line if out of screen
		drawStart := -lineHeight/2 + h/2
		if drawStart < 0 {
			drawStart = 0
		}
		drawEnd := lineHeight/2 + h/2
		if drawEnd >= h {
			drawEnd = h - 1
		}

		line_color := color.RGBA{uint8(210*(1/(dist+1))) + 45, 0, 0, 255}

		vector.StrokeLine(screen, canvas_offset+float32(i), float32(drawStart), canvas_offset+float32(i), float32(drawEnd), 1, line_color, false)

		ray += ray_rad_increment
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	draw2D(g, screen)
	draw3D(g, screen, float32(g.width))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func ebitenGameLoop(file string) {
	body, _ := os.ReadFile(file)
	grid := utils.ParseGrid(string(body))

	scale_factor := 50
	canvas_width := len(grid[0]) * scale_factor
	canvas_height := len(grid) * scale_factor
	window_width := canvas_width * 2
	window_height := canvas_height

	ebiten.SetWindowSize(window_width, window_height)
	ebiten.SetWindowTitle("Raycaster!")

	camera := Camera{
		x:           7,
		y:           7,
		xy_rotation: 0,
		fov:         0.6 * math.Pi,
	}
	game := Game{
		width:                                canvas_width,
		height:                               canvas_height,
		scale_factor:                         scale_factor,
		grid:                                 grid,
		camera:                               camera,
		player_size:                          0.2, // relative to grid unit
		key_down_is_down:                     false,
		key_up_is_down:                       false,
		key_right_is_down:                    false,
		key_left_is_down:                     false,
		key_rotate_clockwise_is_down:         false,
		key_rotate_counter_clockwise_is_down: false,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func main() {
	file := "map.txt"
	ebitenGameLoop(file)
}

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
	// player pos
	x, y float64

	// xy_rotation and fov is in rads
	xy_rotation, fov float64

	// projection plane width and dist from camera, since we only need 2D plane
	plane_width            float64
	plane_dist_from_camera float64
}

type Game struct {
	// screen size and scale factor from grid units to pixels
	width, height, scale_factor int

	// map
	grid [][]bool

	// camera
	camera Camera

	// player diameter
	player_size float64

	// key down states for movement
	key_down_is_down, key_up_is_down, key_left_is_down, key_right_is_down bool
	key_rotate_clockwise_is_down, key_rotate_counter_clockwise_is_down    bool
}

func (g *Game) Update() error {
	// movement
	keyMapping := map[ebiten.Key]*bool{
		// up, down, left, right movement
		ebiten.KeyW:     &g.key_up_is_down,
		ebiten.KeyUp:    &g.key_up_is_down,
		ebiten.KeyS:     &g.key_down_is_down,
		ebiten.KeyDown:  &g.key_down_is_down,
		ebiten.KeyA:     &g.key_left_is_down,
		ebiten.KeyLeft:  &g.key_left_is_down,
		ebiten.KeyD:     &g.key_right_is_down,
		ebiten.KeyRight: &g.key_right_is_down,

		// rotation
		ebiten.KeySpace: &g.key_rotate_clockwise_is_down,
		ebiten.KeyF:     &g.key_rotate_clockwise_is_down,
		ebiten.KeyG:     &g.key_rotate_counter_clockwise_is_down,
	}

	for key, variable := range keyMapping {
		if inpututil.IsKeyJustPressed(key) {
			*variable = true
		}
		if inpututil.IsKeyJustReleased(key) {
			*variable = false
		}
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

	if g.key_rotate_counter_clockwise_is_down {
		g.camera.xy_rotation -= 0.025
	}

	if g.key_rotate_clockwise_is_down {
		g.camera.xy_rotation += 0.025
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
		dist, _, _ := utils.CalculateDistanceOfCastedRay(g.grid, g.camera.x, g.camera.y, ray, false)
		norm_x, norm_y := utils.CalcNormDirVectorFromRadians(ray)
		f_x, f_y := dist*norm_x*float64(g.scale_factor)+g.camera.x*float64(g.scale_factor), dist*norm_y*float64(g.scale_factor)+g.camera.y*float64(g.scale_factor)
		vector.StrokeLine(screen, float32(g.camera.y)*float32(g.scale_factor), float32(g.camera.x)*float32(g.scale_factor), float32(f_y), float32(f_x), 1, color.White, false)
		ray += ray_rad_increment
	}
}

func draw3D(g *Game, screen *ebiten.Image, canvas_offset float32) {
	// start_ray_rad := g.camera.xy_rotation - (g.camera.fov / 2)
	// ray_rad_increment := g.camera.fov / (float64(g.width) - 1)
	h := g.height
	// ray := start_ray_rad

	target_plane_x := 0.5
	for i := 0; i < g.width; i++ {
		// we will calculate the radians of the ray that will go through this pixel
		// in the camera plan
		target_x_from_plane_center := target_plane_x - g.camera.plane_width/2
		radian_offset_from_xy_rotation := math.Atan(target_x_from_plane_center / g.camera.plane_dist_from_camera)
		ray := radian_offset_from_xy_rotation + g.camera.xy_rotation

		// the euclidean dist has some distortion
		euclidean_dist, _, _ := utils.CalculateDistanceOfCastedRay(g.grid, g.camera.x, g.camera.y, ray, false)
		// lets correct the fisheye distortion using some trignometry to get the perpendicular height
		dist := euclidean_dist * math.Cos(radian_offset_from_xy_rotation)

		// Height of line to draw
		lineHeight := int(float64(h) / dist / 2)

		// Find start and end of line, cap the line if out of screen
		drawStart := -lineHeight/2 + h/2
		if drawStart < 0 {
			drawStart = 0
		}
		drawEnd := lineHeight/2 + h/2
		if drawEnd >= h {
			drawEnd = h - 1
		}

		line_color := color.RGBA{uint8(215*(1/(dist+1))) + 40, 0, 0, 255}

		vector.StrokeLine(screen, canvas_offset+float32(i), float32(drawStart), canvas_offset+float32(i), float32(drawEnd), 1, line_color, false)

		// ray += ray_rad_increment
		target_plane_x += 1
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	draw2D(g, screen)
	draw3D(g, screen, float32(g.width))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func ebitenGameLoop(grid [][]bool) {
	scale_factor := 50
	canvas_width := len(grid[0]) * scale_factor
	canvas_height := len(grid) * scale_factor
	window_width := canvas_width * 2
	window_height := canvas_height

	ebiten.SetWindowSize(window_width, window_height)
	ebiten.SetWindowTitle("Raycaster!")

	fov := 0.5 * math.Pi
	plane_width := float64(canvas_width)
	plane_dist_from_camera := plane_width / (2 * math.Tan(fov/2))

	camera := Camera{
		x:                      6,
		y:                      6,
		xy_rotation:            0,
		fov:                    fov,
		plane_width:            float64(canvas_width),
		plane_dist_from_camera: plane_dist_from_camera,
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
	body, _ := os.ReadFile(file)
	world := utils.ParseGrid(string(body))

	ebitenGameLoop(world)
}

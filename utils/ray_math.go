package utils

import (
	"fmt"
	"math"
)

// will turn the xy_rotation, which is 0 when pointing straight towards N,
// into dir vector on this coordinate plane
func CalcNormDirVectorFromRadians(xy_rotation float64) (float64, float64) {
	// this is a linear transformation over the x axis
	return -1 * math.Cos(xy_rotation), math.Sin(xy_rotation)
}

// Get the perpendicular distance from the camera
// useful to correct fish eye effect
func GetPerpendicularDistance(distance, ray_angle, fov float64) float64 {
	// TODO: Implement fish eye correction
	return distance * math.Cos(ray_angle)
}

func GetTargetInStepDirection(x, y float64, x_step, y_step int) (float64, float64) {
	target_x, target_y := x, y

	// handle floating point errors
	if math.Abs(x-math.Round(x)) <= 0.001 {
		x = math.Round(x)
	}
	if math.Abs(y-math.Round(y)) <= 0.001 {
		y = math.Round(y)
	}

	if x_step < 0 {
		target_x = math.Ceil(x - 1)
	} else if x_step > 0 {
		target_x = math.Floor(x + 1)
	}

	if y_step < 0 {
		target_y = math.Ceil(y - 1)
	} else if y_step > 0 {
		target_y = math.Floor(y + 1)
	}

	return target_x, target_y
}

func CanContinueInDirection(grid [][]bool, x, y float64, ray_dir_x, ray_dir_y float64) bool {
	// walk a tiny step more in the ray direction and check if cell is wall
	step := 0.05
	x += ray_dir_x * step
	y += ray_dir_y * step

	if x < 0 || y < 0 || x >= float64(len(grid)) || y >= float64(len(grid[0])) || grid[int(math.Floor(x))][int(math.Floor(y))] {
		return false
	}

	return true
}

func debugPrint(print bool, a ...any) {
	if print {
		fmt.Println(a...)
	}
}

func DDA(grid [][]bool, x, y, xy_rotation float64, debugging bool) (float64, float64, float64) {
	// need to calculate the norm vector direction from the xy_rotation
	ray_dir_x, ray_dir_y := CalcNormDirVectorFromRadians(xy_rotation)

	// how much along the line we go if we go 1 unit in x
	x_unit_step := math.Sqrt(1 + (ray_dir_y/ray_dir_x)*(ray_dir_y/ray_dir_x))

	// how much along the line we go if we go 1 unit in y
	y_unit_step := math.Sqrt(1 + (ray_dir_x/ray_dir_y)*(ray_dir_x/ray_dir_y))

	// our current position
	cur_x, cur_y := x, y

	// which way we are stepping
	x_step, y_step := 0, 0

	// set the step
	if ray_dir_x < 0 {
		x_step = -1
	} else if ray_dir_x > 0 {
		x_step = 1
	}

	if ray_dir_y < 0 {
		y_step = -1
	} else if ray_dir_y > 0 {
		y_step = 1
	}

	wall_tile_found := false
	iter := 0
	max_iter := len(grid) * len(grid[0])
	distance := 0.0

	for !wall_tile_found && iter < max_iter {
		iter += 1

		t_x, t_y := GetTargetInStepDirection(cur_x, cur_y, x_step, y_step)
		dist_x, dist_y := math.Abs(t_x-cur_x), math.Abs(t_y-cur_y)
		line_x, line_y := dist_x*x_unit_step, dist_y*y_unit_step

		debugPrint(debugging, "At:", cur_x, cur_y, "\t Towards:", t_x, t_y, "\t X Option:", dist_x, line_x, "\t Y Option:", dist_y, line_y)

		if line_x < line_y {
			debugPrint(debugging, "Choosing X", cur_x, cur_y, t_x, t_y)
			distance += line_x
			cur_x = t_x
			cur_y = y + float64(ray_dir_y)*(float64(t_x-x)/ray_dir_x)
			debugPrint(debugging, "Choose X", cur_x, cur_y)
		} else {
			debugPrint(debugging, "Choosing Y", cur_x, cur_y, t_x, t_y)
			distance += line_y
			cur_y = t_y
			cur_x = x + float64(ray_dir_x)*(float64(t_y-y)/ray_dir_y)
			debugPrint(debugging, "Choose Y", cur_x, cur_y)
		}

		if !CanContinueInDirection(grid, cur_x, cur_y, ray_dir_x, ray_dir_y) {
			wall_tile_found = true
		}
	}

	return distance, cur_x, cur_y
}

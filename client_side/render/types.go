package render

import (
	"github.com/Jla3eP/tetris/both_sides_code"
)

type (
	Render struct {
		FieldSize both_sides_code.Coords2
	}

	renderConfig struct {
		TexturePackName           string `json:"texturepack_name"`
		TextureWidth              int32  `json:"texture_width"`
		TextureHeight             int32  `json:"texture_height"`
		BonusPxWindowWidth        int32  `json:"bonus_px_window_width"`
		BonusPxWindowHeight       int32  `json:"bonus_px_window_height"`
		BonusPercentsWindowWidth  int32  `json:"bonus_percents_window_width"`
		BonusPercentsWindowHeight int32  `json:"bonus_percents_window_height"`
		FieldBiasX                int32  `json:"field_bias_x"`
		FieldBiasY                int32  `json:"field_bias_y"`
		TargetFps                 int32  `json:"target_fps"`
		PrintFps                  bool   `json:"print_fps"`
		PixelsToEnemiesField      int32  `json:"pixels_to_enemies_field"`
	}
)

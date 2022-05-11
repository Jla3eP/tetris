package render

import (
	"github.com/Jla3eP/tetris/client_side/field"
)

type Render struct {
	FieldSize field.Coords2
}

type renderConfig struct {
	TexturePackName           string `json:"texturepack_name"`
	TextureWidth              int32  `json:"texture_width"`
	TextureHeight             int32  `json:"texture_height"`
	BonusPxWindowWidth        int32  `json:"bonus_px_window_width"`
	BonusPxWindowHeight       int32  `json:"bonus_px_window_height"`
	BonusPercentsWindowWidth  int32  `json:"bonus_percents_window_width"`
	BonusPercentsWindowHeight int32  `json:"bonus_percents_window_height"`
	FieldBiasX                uint32 `json:"field_bias_x"`
	FieldBiasY                uint32 `json:"field_bias_y"`
	TargetFps                 int32  `json:"target_fps"`
	PrintFps                  bool   `json:"print_fps"`
}

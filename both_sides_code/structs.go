package both_sides_code

import "sync"

type (
	Coords2 struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	AuthInfo struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
	}

	FieldRequest struct {
		EnemyFigureID          int     `json:"enemy_figure_id"`
		EnemyFigureColor       int     `json:"enemy_figure_color"`
		EnemyFigureRotateIndex int     `json:"enemy_figure_rotate_index"`
		EnemyFigureCoords      Coords2 `json:"enemy_figure_coords"`
		EnemyFigureSent        bool    `json:"enemy_figure_sent"`

		Nickname   string `json:"nickname"`
		SessionKey string `json:"session_key"`
	}

	FieldResponse struct {
		FigureID    int `json:"figure_id"`
		FigureColor int `json:"figure_color"`

		FieldRequest `json:"field_request"`
	}

	ResponseStatus struct {
		Comment string `json:"comment"`
	}

	SessionUpdateRequest struct {
		Nickname   string `json:"nickname"`
		SessionKey string `json:"session_key"`
	}

	PossibleStatus struct {
		Coords []Coords2 `json:"vec2"`
	}

	Figure struct {
		ID                 int8
		PossibleStatuses   []PossibleStatus `json:"possible_statuses"`
		Color              int
		CurrentRotateIndex int
		CurrentCoords      Coords2
		Mutex              *sync.Mutex
		Fixed              bool
	}

	FiguresConfig struct {
		Figures []Figure `json:"figures"`
	}
)

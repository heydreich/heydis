package sortedset

import (
	"errors"
	"strconv"
)

const (
	negativeInf int8 = -1
	positiveInf int8 = 1
)

// ScoreBorder represents range of a float value, including: <, <=, >, >=, +inf, -inf
type ScoreBorder struct {
	Inf     int8
	Value   float64
	Exclude bool
}

func (border *ScoreBorder) greater(value float64) bool {
	if border.Inf == negativeInf {
		return false
	} else if border.Inf == positiveInf {
		return true
	}
	if border.Exclude {
		return border.Value > value
	}
	return border.Value >= value
}

func (border *ScoreBorder) less(value float64) bool {
	if border.Inf == positiveInf {
		return false
	} else if border.Inf == negativeInf {
		return true
	}
	if border.Exclude {
		return border.Value < value
	}
	return border.Value <= value
}

var positiveInfBorder = &ScoreBorder{
	Inf: positiveInf,
}
var negativeInfBorder = &ScoreBorder{
	Inf: negativeInf,
}

// ParseScoreBorder creates ScoreBorder from redis arguments
func ParseScoreBorder(s string) (*ScoreBorder, error) {
	if s == "inf" || s == "+inf" {
		return positiveInfBorder, nil
	}
	if s == "-inf" {
		return negativeInfBorder, nil
	}
	if s[0] == '(' {
		value, err := strconv.ParseFloat(s[1:], 64)
		if err != nil {
			return nil, errors.New("ERR min or max is not a float")
		}
		return &ScoreBorder{
			Inf:     0,
			Value:   value,
			Exclude: true,
		}, nil
	}
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.New("ERR min or max is not a float")
	}
	return &ScoreBorder{
		Inf:     0,
		Value:   value,
		Exclude: false,
	}, nil
}

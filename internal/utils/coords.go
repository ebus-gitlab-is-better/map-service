package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mojixcoder/gosrm"
)

func ParseCoordinates(input string) ([]gosrm.Coordinate, error) {
	var points []gosrm.Coordinate
	coords := strings.Split(input, ";")
	for _, coord := range coords {
		parts := strings.Split(coord, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid coordinate format")
		}
		lon, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, err
		}
		lat, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, err
		}
		points = append(points, gosrm.Coordinate{lon, lat})
	}
	return points, nil
}
func distanceToSegment(p, a, b [2]float64) float64 {
	ap := [2]float64{p[0] - a[0], p[1] - a[1]}
	ab := [2]float64{b[0] - a[0], b[1] - a[1]}
	ab2 := ab[0]*ab[0] + ab[1]*ab[1]
	ap_ab := ap[0]*ab[0] + ap[1]*ab[1]
	t := ap_ab / ab2
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	nearest := [2]float64{a[0] + ab[0]*t, a[1] + ab[1]*t}
	return math.Sqrt((nearest[0]-p[0])*(nearest[0]-p[0]) + (nearest[1]-p[1])*(nearest[1]-p[1]))
}

// Проверка, находится ли точка близко к линии.
func IsPointNearPolyline(point [2]float64, shape string, tolerance float64) (bool, error) {
	coords := DecodePolyline(&shape)
	// fmt.Print(coords)

	for i := 0; i < len(coords)-1; i++ {
		if distanceToSegment(point, coords[0], coords[i+1]) < tolerance {
			return true, nil
		}
	}

	return false, nil
}
func DecodePolyline(encoded *string, precisionOptional ...int) [][2]float64 {
	// default to 6 digits of precision
	precision := 6
	if len(precisionOptional) > 0 {
		precision = precisionOptional[0]
	}
	factor := math.Pow10(precision)

	// Coordinates have variable length when encoded, so just keep
	// track of whether we've hit the end of the string. In each
	// loop iteration, a single coordinate is decoded.
	lat, lng := 0, 0
	var coordinates [][2]float64
	index := 0
	for index < len(*encoded) {
		// Consume varint bits for lat until we run out
		var byte int = 0x20
		shift, result := 0, 0
		for byte >= 0x20 {
			byte = int((*encoded)[index]) - 63
			result |= (byte & 0x1f) << shift
			shift += 5
			index++
		}

		// check if we need to go negative or not
		if (result & 1) > 0 {
			lat += ^(result >> 1)
		} else {
			lat += result >> 1
		}

		// Consume varint bits for lng until we run out
		byte = 0x20
		shift, result = 0, 0
		for byte >= 0x20 {
			byte = int((*encoded)[index]) - 63
			result |= (byte & 0x1f) << shift
			shift += 5
			index++
		}

		// check if we need to go negative or not
		if (result & 1) > 0 {
			lng += ^(result >> 1)
		} else {
			lng += result >> 1
		}

		// scale the int back to floating point and store it
		coordinates = append(coordinates, [2]float64{float64(lat) / factor, float64(lng) / factor})
	}

	return coordinates
}

func EncodePolyline(coordinates [][2]float64, precisionOptional ...int) string {
	precision := 6
	if len(precisionOptional) > 0 {
		precision = precisionOptional[0]
	}
	factor := math.Pow10(precision)

	var encoded strings.Builder
	previousLat, previousLng := 0, 0

	for _, coordinate := range coordinates {
		lat := int(math.Round(coordinate[0] * factor))
		lng := int(math.Round(coordinate[1] * factor))
		encodeValue(&encoded, lat-previousLat)
		encodeValue(&encoded, lng-previousLng)
		previousLat, previousLng = lat, lng
	}

	return encoded.String()
}

func encodeValue(encoded *strings.Builder, value int) {
	value = value << 1
	if value < 0 {
		value = ^value
	}
	for value >= 0x20 {
		encoded.WriteByte(byte((0x20 | (value & 0x1f)) + 63))
		value >>= 5
	}
	encoded.WriteByte(byte(value + 63))
}

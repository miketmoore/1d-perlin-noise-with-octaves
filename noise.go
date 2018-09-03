package noise

func GenerateNoise(amp, wl, octaves, divisor, width float64) []Perlin {
	result := []Perlin{}
	for i := 0.0; i < octaves; i++ {
		result = append(result, NewPerlin(amp, wl, width))
		amp = amp / divisor
		wl = wl / divisor
	}
	return result
}

func CombineNoise(pl []Perlin) []float64 {
	result := []float64{}
	for i, total, j := 0, 0.0, 0; i < len(pl[0].Pos); i++ {
		total = 0
		for j = 0; j < len(pl); j++ {
			total += pl[j].Pos[i]
		}
		result = append(result, total)
	}
	return result
}

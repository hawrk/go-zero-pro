package util

// CalOptimizeScore  计算优选算法的数值 todo 该计算是临时的，待修改
func CalOptimizeScore(openRate float64, incomeRate float64, basisPoint float64) float64 {
	return (openRate + incomeRate + basisPoint) / 3
}

package textprocessor

import (
	"regexp"
	"sort"
	"strings"
)

// StopWords 是一个常见的停用词集合
var StopWords = map[string]bool{
	"the": true, "and": true, "is": true, "in": true,
	"on": true, "at": true, "to": true, "of": true,
}

// WordCount 表示一个单词及其出现次数
type WordCount struct {
	Word  string
	Count int
}

// ProcessText 处理输入文本，提取有效词汇，统计词频并按频率排序
func ProcessText(input string) []WordCount {
	// 清洗文本：转小写 + 去除非字母字符
	reg, _ := regexp.Compile("[^a-zA-Z\\s]+")
	cleaned := reg.ReplaceAllString(input, "")
	lower := strings.ToLower(cleaned)

	// 分割为单词
	words := strings.Fields(lower)

	// 统计词频
	wordCounts := make(map[string]int)
	for _, word := range words {
		if StopWords[word] {
			continue
		}
		wordCounts[word]++
	}

	// 转换为切片用于排序
	var result []WordCount
	for word, count := range wordCounts {
		result = append(result, WordCount{Word: word, Count: count})
	}

	// 按照词频降序排序，如果频率相同则按字母升序
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Word < result[j].Word
		}
		return result[i].Count > result[j].Count
	})

	return result
}

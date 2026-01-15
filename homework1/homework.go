package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func main() {
	nums := []int{1, 3, 4, 5, 3, 2, 4, 5, 2, 1, 9}
	fmt.Printf("single number is: %d\n", SingleNumber(nums))
}

// 1. 只出现一次的数字
// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
func SingleNumber(nums []int) int {
	m := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		if _, ok := m[nums[i]]; ok {
			delete(m, nums[i])

		} else {
			m[nums[i]] = 1
		}

	}
	var num int
	for k := range m {
		num = k
		break // 拿到第一个后立即跳出
	}
	return num
}

// 2. 回文数
// 判断一个整数是否是回文数
func IsPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	s := strconv.Itoa(x)
	lenth := len(s)

	for i := 0; i < lenth/2; i++ {
		if s[i] != s[lenth-1-i] {
			return false
		}
	}
	return true
}

// 3. 有效的括号
// 给定一个只包括 '(', ')', '{', '}', '[', ']' 的字符串，判断字符串是否有效
func IsValid(s string) bool {
	lenth := len(s)
	if lenth%2 != 0 {
		return false
	}
	bracketMap := map[byte]byte{
		'(': ')',
		'{': '}',
		'[': ']',
	}
	stack := []byte{}

	for i := 0; i < lenth; i++ {
		if _, ok := bracketMap[s[i]]; ok {
			//encounter open, add to stack
			stack = append(stack, s[i])
		} else {
			//encounter close, get top element of stack
			topElement := stack[len(stack)-1]
			if s[i] != bracketMap[topElement] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return true
}

// 4. 最长公共前缀
// 查找字符串数组中的最长公共前缀
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := ""

	for j := 1; j <= len(strs[0]); j++ {
		prefix = strs[0][:j]

		for i := 0; i <= len(strs)-1; i++ {
			if !strings.HasPrefix(strs[i], prefix) {
				return prefix[:j-1]
			}
		}
	}
	return prefix
}

// 5. 加一
// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func PlusOne(digits []int) []int {
	result := 0
	for _, num := range digits {
		result = result*10 + num
	}
	result = result + 1

	var resultArr []int
	for result > 0 {
		resultArr = append(resultArr, result%10)
		result = result / 10
	}
	for i, j := 0, len(resultArr)-1; i < j; i, j = i+1, j-1 {
		resultArr[i], resultArr[j] = resultArr[j], resultArr[i]
	}
	return resultArr
}

// 6. 删除有序数组中的重复项 0 1 2
// 给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
func RemoveDuplicates(nums []int) int {

	n := len(nums)
	if n == 0 {
		return 0
	}

	i := 0
	for j := 0; j < len(nums); j++ {
		if nums[i] != nums[j] {
			i++
			nums[i] = nums[j]
		}
	}
	return i + 1

}

// 7. 合并区间
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
func Merge(intervals [][]int) [][]int {
	if len(intervals) <= 1 {
		return intervals
	}
	slices.SortFunc(intervals, func(a, b []int) int {
		return a[0] - b[0]
	})

	// 2. 初始化结果集，先放入第一个区间
	res := [][]int{intervals[0]}

	for i := 1; i < len(intervals); i++ {
		curr := intervals[i]
		// 获取结果集中最后一个区间的引用
		last := &res[len(res)-1]

		// 3. 判断是否重叠：当前区间的起点 <= 结果集中最后一个区间的终点
		if curr[0] <= (*last)[1] {
			// 如果重叠，更新结果集中最后一个区间的终点为两者的最大值
			if curr[1] > (*last)[1] {
				(*last)[1] = curr[1]
			}
		} else {
			// 4. 如果不重叠，直接将当前区间加入结果集
			res = append(res, curr)
		}
	}
	return res
}

// 8. 两数之和
// 给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
func TwoSum(nums []int, target int) []int {
	// TODO: implement
	return nil
}

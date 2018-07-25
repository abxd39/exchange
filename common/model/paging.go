package model

import "math"

// 分页结构
type ModelList struct {
	IsPage    bool        `json:"is_page"`    // 是否分页
	PageIndex int         `json:"page_index"` // 当前页码
	PageSize  int         `json:"page_size"`  // 每页数据条数
	PageCount int         `json:"page_count"` // 总页数
	Total     int         `json:"total"`      // 总数据条数
	Items     interface{} `json:"items"`      // 数据数组
}

// 分页列表
func Paging(pageIndex, pageSize, total int) (offset int, modelList *ModelList) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset = (pageIndex - 1) * pageSize

	modelList = &ModelList{
		IsPage:    true,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		PageCount: int(math.Ceil(float64(total) / float64(pageSize))),
		Total:     total,
	}

	return
}

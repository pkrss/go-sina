
package beans


type Page struct {
	Content				interface{}			`json:"content"`	// 内容
	Last				bool				`json:"last"`		// 是否是最后一页
	TotalPages			int					`json:"totalPages"`	// 总页数
	TotalElements		int					`json:"totalElements"`	// 记录总条数
	Sort				[]interface{}		`json:"sort"`	// 排序参数
	First				bool				`json:"first"`			// 是否是第一页
	NumberOfElements 	int					`json:"numberOfElements"`		// 返回的content数组大小
	Size				int					`json:"size"`		// 限制大小
	Number				int					`json:"number"`		// 起始页，0开始
}
package errcode

var (
	ErrorGetObjectListFail = NewError(20010001, "获取对象列表失败")
	ErrorCreateObjectFail  = NewError(20010002, "创建对象失败")
	ErrorUpdateObjectFail  = NewError(20010003, "更新对象失败")
	ErrorDeleteObjectFail  = NewError(20010004, "删除对象失败")
	ErrorCountObjectFail   = NewError(20010005, "统计对象失败")
	ErrorUnShortObjectFail = NewError(20010006, "还原对象失败")
	ErrorGetObjectFail     = NewError(20010007, "获取对象失败")
	ErrorPutObjectFail     = NewError(20010008, "上传对象失败")

	ErrorUploadFileFail = NewError(20030001, "上传文件失败")
)

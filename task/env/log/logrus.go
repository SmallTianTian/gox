package log

// LogrusTask 新增 logrus 日志.
// 需要前置检查：
// 1. 该地址是否是 go 项目
// 2. 该项目存在 pkg/logger 目录
// 3. 该项目中不存在 logrus 日志文件
//
// 主要做以下几件事情：
// 1. 新建 logrus 日志文件
// 2. 如果不存在 logoption 文件，则新建。
// 3. 在 main 文件中，写入该日志初始化
// 4. 更新 config
type LogrusTask struct {
	commonTask
}

func NewLogrusTask() *LogrusTask {
	return &LogrusTask{
		commonTask: *newCommonTask("logrus"),
	}
}

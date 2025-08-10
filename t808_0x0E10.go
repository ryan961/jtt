package jtt

// T808_0x0E10 驾驶员身份识别上报
type T808_0x0E10 struct {
	// 比对结果
	// 0：匹配成功
	// 1：匹配失败
	// 2：超时
	// 3：没有启用该功能
	// 4：连接异常
	// 5：无指定人脸图片
	// 6：无人脸库
	MatchResult byte `json:"matchResult"`
	// 比对相似度阈值
	// 百分比；范围 0%～100%。单位是 1%
	ThresholdPercentage byte `json:"thresholdPercentage"`
	// 比对相似度
	// 百分比；范围 0.00%～100.00%。单位是 0.01%；比如 5432 表示 54.32%
	SimilarityPercentage uint16 `json:"similarityPercentage"`
	// 比对类型
	// 0-插卡比对
	// 1-巡检比对
	// 2-点火比对
	// 3-离开返回比对
	MatchType byte `json:"matchType"`
	// 比对人脸 ID
	MatchFaceID string `json:"matchFaceID"`
	// 位置信息汇报（0x0200）消息体
	// 表示人脸比对时刻的位置基本信息数据
	LocationInfo T808_0x0200 `json:"locationInfo"`
	// 图片格式
	// 0: JPEG
	ImageFormat byte `json:"imageFormat"`
	// 图片数据包
	// 比对结果为 0 或者 1 时，应上传图片数据（为抓拍的图片）
	ImageData []byte `json:"imageData"`
}

func (msg *T808_0x0E10) MsgID() MsgID {
	return MsgT808_0x0E10
}

func (msg *T808_0x0E10) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入比对结果
	writer.WriteByte(msg.MatchResult)

	// 写入比对相似度阈值
	writer.WriteByte(msg.ThresholdPercentage)

	// 写入比对相似度
	writer.WriteWord(msg.SimilarityPercentage)

	// 写入比对类型
	writer.WriteByte(msg.MatchType)

	// 写入比对人脸 ID 长度
	writer.WriteByte(byte(len(msg.MatchFaceID)))

	// 写入比对人脸 ID
	err := writer.WriteString(msg.MatchFaceID)
	if err != nil {
		return nil, err
	}

	// 写入位置信息
	locationData, err := msg.LocationInfo.Encode()
	if err != nil {
		return nil, err
	}
	writer.Write(locationData)

	// 写入图片格式
	writer.WriteByte(msg.ImageFormat)

	// 写入图片数据
	writer.Write(msg.ImageData)

	return writer.Bytes(), nil
}

func (msg *T808_0x0E10) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	// 读取比对结果
	var err error
	msg.MatchResult, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取比对相似度阈值
	msg.ThresholdPercentage, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取比对相似度
	msg.SimilarityPercentage, err = reader.ReadWord()
	if err != nil {
		return 0, err
	}

	// 读取比对类型
	msg.MatchType, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取比对人脸 ID 长度
	faceIDLen, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取比对人脸 ID
	msg.MatchFaceID, err = reader.ReadString(int(faceIDLen))
	if err != nil {
		return 0, err
	}

	// 读取位置信息
	locationData, err := reader.Read(28) // 位置信息固定长度为28字节
	if err != nil {
		return 0, err
	}

	msg.LocationInfo = T808_0x0200{}
	_, err = msg.LocationInfo.Decode(locationData)
	if err != nil {
		return 0, err
	}

	// 读取图片格式
	msg.ImageFormat, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取图片数据
	if reader.Len() > 0 {
		msg.ImageData, err = reader.Read(reader.Len())
		if err != nil {
			return 0, err
		}
	}

	return len(data) - reader.Len(), nil
}

package jtt

// T808_0x8E11 驾驶员身份信息下发
type T808_0x8E11 struct {
	// 设置类型
	// 0：增加（全替换）
	// 1：删除（全删除）
	// 2：删除指定条目
	// 3：修改（如果设备存在人脸 id，那么替换当前人脸特征码。如果设备不存在人脸 id，那么新增人脸）
	SetType byte `json:"setType"`
	// 驾驶员库列表个数
	DriverCount byte `json:"driverCount"`
	// 驾驶员库信息列表
	DriverInfoList []T808_0x8E11_DriverInfo `json:"driverInfoList"`
}

// T808_0x8E11_DriverInfo 人脸信息列表数据格式
type T808_0x8E11_DriverInfo struct {
	// 人脸 ID
	FaceID string `json:"faceID"`
	// 从业资格证
	DriverLicense string `json:"driverLicense"`
	// 人脸图片地址协议
	// 0--FTP, 1--HTTP
	FaceImageProtocol byte `json:"faceImageProtocol"`
	// 人脸图片地址
	FaceImageURL string `json:"faceImageURL"`
	// 人脸图片来源
	// 0--本机拍摄图片，1--第三方图片
	FaceImageSource byte `json:"faceImageSource"`
}

func (msg *T808_0x8E11) MsgID() MsgID {
	return MsgT808_0x8E11
}

func (msg *T808_0x8E11) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入设置类型
	writer.WriteByte(msg.SetType)

	// 写入驾驶员库列表个数
	writer.WriteByte(msg.DriverCount)

	// 写入驾驶员库信息列表
	for _, driverInfo := range msg.DriverInfoList {
		// 写入人脸 ID 长度
		writer.WriteByte(byte(len(driverInfo.FaceID)))
		// 写入人脸 ID
		err := writer.WriteString(driverInfo.FaceID)
		if err != nil {
			return nil, err
		}

		// 写入从业资格证长度
		writer.WriteByte(byte(len(driverInfo.DriverLicense)))
		// 写入从业资格证
		err = writer.WriteString(driverInfo.DriverLicense)
		if err != nil {
			return nil, err
		}

		// 写入人脸图片地址协议
		writer.WriteByte(driverInfo.FaceImageProtocol)

		// 写入人脸图片地址长度
		writer.WriteByte(byte(len(driverInfo.FaceImageURL)))
		// 写入人脸图片地址
		err = writer.WriteString(driverInfo.FaceImageURL)
		if err != nil {
			return nil, err
		}

		// 写入人脸图片来源
		writer.WriteByte(driverInfo.FaceImageSource)
	}

	return writer.Bytes(), nil
}

func (msg *T808_0x8E11) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	// 读取设置类型
	var err error
	msg.SetType, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取驾驶员库列表个数
	msg.DriverCount, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取驾驶员库信息列表
	msg.DriverInfoList = make([]T808_0x8E11_DriverInfo, msg.DriverCount)
	for i := 0; i < int(msg.DriverCount); i++ {
		// 读取人脸 ID 长度
		faceIDLen, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取人脸 ID
		msg.DriverInfoList[i].FaceID, err = reader.ReadString(int(faceIDLen))
		if err != nil {
			return 0, err
		}

		// 读取从业资格证长度
		driverLicenseLen, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取从业资格证
		msg.DriverInfoList[i].DriverLicense, err = reader.ReadString(int(driverLicenseLen))
		if err != nil {
			return 0, err
		}

		// 读取人脸图片地址协议
		msg.DriverInfoList[i].FaceImageProtocol, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取人脸图片地址长度
		faceImageURLLen, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取人脸图片地址
		msg.DriverInfoList[i].FaceImageURL, err = reader.ReadString(int(faceImageURLLen))
		if err != nil {
			return 0, err
		}

		// 读取人脸图片来源
		msg.DriverInfoList[i].FaceImageSource, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}
	}

	return len(data) - reader.Len(), nil
}

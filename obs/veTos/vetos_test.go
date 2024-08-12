package veTos

//
//func UploadFileT(c *gin.Context) {
//	// 获取表单中的文件
//	file, err := c.FormFile("img")
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{"code": common.PARAM_ERROR, "msg": "failed to get form"})
//		return
//	}
//	// 保存文件到本地临时路径
//	dst := filepath.Join("./obsData", file.Filename)
//	if err := c.SaveUploadedFile(file, dst); err != nil {
//		c.JSON(http.StatusOK, gin.H{"code": common.PARAM_ERROR, "msg": "failed to save local file"})
//		return
//	}
//	// 调用 veTos.UploadTos 上传文件
//	if err := veTos.UploadTos(dst); err != nil {
//		c.JSON(http.StatusOK, gin.H{"code": common.PARAM_ERROR, "msg": "failed to upload tos"})
//		return
//	}
//	// 上传成功后删除本地文件
//	if err := os.Remove(dst); err != nil {
//		c.JSON(http.StatusOK, gin.H{"code": common.PARAM_ERROR, "msg": "failed to delete file"})
//		return
//	}
//	//下载文件
//	//veTos.DownloadTos()
//
//	c.JSON(http.StatusOK, gin.H{"code": common.SUCCESSFUL, "msg": "ok"})
//}

package postprocess

func GetErrorDetails(moduleType string) string {
	switch moduleType {
	case string(ModulePromotion):
		return "促销数据异常！"
	case string(ModuleMileage):
		return "积分数据异常！"
	case string(ModulePay):
		return "支付数据异常！"
	case string(ModuleSalePerson):
		return "营业员销售业绩异常！"
	case string(ModuleSaleFee):
		return "扣率数据异常！"
	case string(ModuleRefundApproval):
		return "退货审批数据异常！"
	case string(SendToClearance):
		return "上传数据异常！"
	}
	return "未知异常！"
}

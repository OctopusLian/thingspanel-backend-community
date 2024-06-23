package service

import (
	"project/dal"
	"project/model"
	"project/others/http_client"
	"project/query"
	utils "project/utils"
	"strconv"
	"time"

	"github.com/go-basic/uuid"
	"github.com/jinzhu/copier"
)

type ServiceAccess struct{}

func (s *ServiceAccess) CreateAccess(req *model.CreateAccessReq) (map[string]interface{}, error) {
	var serviceAccess model.ServiceAccess
	copier.Copy(&serviceAccess, req)
	serviceAccess.ID = uuid.New()
	if *serviceAccess.ServiceAccessConfig == "" {
		*serviceAccess.ServiceAccessConfig = "{}"
	}
	serviceAccess.CreateAt = time.Now().UTC()
	serviceAccess.UpdateAt = time.Now().UTC()
	err := query.ServiceAccess.Create(&serviceAccess)
	if err != nil {
		return nil, err
	}
	resp := make(map[string]interface{})
	resp["id"] = serviceAccess.ID
	return resp, nil
}

func (s *ServiceAccess) List(req *model.GetServiceAccessByPageReq) (map[string]interface{}, error) {
	total, list, err := dal.GetServiceAccessListByPage(req)
	listRsp := make(map[string]interface{})
	listRsp["total"] = total
	listRsp["list"] = list

	return listRsp, err
}

func (s *ServiceAccess) Update(req *model.UpdateAccessReq) error {
	updates := make(map[string]interface{})
	updates["service_access_config"] = req.ServiceAccessConfig
	updates["update_at"] = time.Now().UTC()
	err := dal.UpdateServiceAccess(req.ID, updates)
	return err
}

func (s *ServiceAccess) Delete(req *model.DeleteAccessReq) error {
	err := dal.DeleteServiceAccess(req.ID)
	return err
}

// GetVoucherForm
func (s *ServiceAccess) GetVoucherForm(req *model.GetServiceAccessVoucherFormReq) (interface{}, error) {
	// 根据service_plugin_id获取插件服务信息http地址
	servicePlugin, httpAddress, err := dal.GetServicePluginHttpAddressByID(req.ServicePluginID)
	if err != nil {
		return nil, err
	}
	return http_client.GetPluginFromConfigV2(httpAddress, servicePlugin.ServiceIdentifier, "", "SVCRT")
}

// GetServiceAccessDeviceList
func (s *ServiceAccess) GetServiceAccessDeviceList(req *model.ServiceAccessDeviceListReq, userClaims *utils.UserClaims) (interface{}, error) {
	// 通过voucher获取service_plugin_id
	serviceAccess, err := dal.GetServiceAccessByVoucher(req.Voucher, userClaims.TenantID)
	if err != nil {
		return nil, err
	}
	// 根据service_plugin_id获取插件服务信息的http地址
	_, httpAddress, err := dal.GetServicePluginHttpAddressByID(serviceAccess.ID)
	if err != nil {
		return nil, err
	}
	// 数字转字符串
	return http_client.GetServiceAccessDeviceList(httpAddress, req.Voucher, strconv.Itoa(req.PageSize), strconv.Itoa(req.Page))
}
